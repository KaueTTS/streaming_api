package security_test

import (
	"testing"
	"time"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	security "github.com/KaueTTS/streaming_api/src/security"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	t.Run("should generate valid token", func(t *testing.T) {
		setupJWTEnv(t)

		token, err := security.GenerateToken(1, "user@email.com", "user")

		require.NoError(t, err)
		assert.NotEmpty(t, token)

		claims, err := security.ValidateToken(token)

		require.NoError(t, err)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, "user@email.com", claims.Email)
		assert.Equal(t, "user", claims.Role)
		assert.Equal(t, "streaming_api_test", claims.Issuer)
		assert.NotNil(t, claims.IssuedAt)
		assert.NotNil(t, claims.ExpiresAt)
		assert.WithinDuration(t, time.Now().Add(2*time.Hour), claims.ExpiresAt.Time, time.Second)
	})

	t.Run("should return error when jwt secret is empty", func(t *testing.T) {
		setupJWTEnv(t)
		env.JWTSecret = " "

		token, err := security.GenerateToken(1, "user@email.com", "user")

		assert.Empty(t, token)
		assert.EqualError(t, err, "JWT_SECRET n\u00e3o informado")
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("should return claims when token is valid", func(t *testing.T) {
		setupJWTEnv(t)

		token := makeToken(t, jwt.SigningMethodHS256, security.Claims{
			UserID: 10,
			Email:  "admin@email.com",
			Role:   "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    env.AppName,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		claims, err := security.ValidateToken(token)

		require.NoError(t, err)
		assert.Equal(t, uint(10), claims.UserID)
		assert.Equal(t, "admin@email.com", claims.Email)
		assert.Equal(t, "admin", claims.Role)
	})

	t.Run("should return error when jwt secret is empty", func(t *testing.T) {
		setupJWTEnv(t)
		env.JWTSecret = ""

		claims, err := security.ValidateToken("token")

		assert.Nil(t, claims)
		assert.EqualError(t, err, "JWT_SECRET n\u00e3o informado")
	})

	t.Run("should return error when token is expired", func(t *testing.T) {
		setupJWTEnv(t)

		token := makeToken(t, jwt.SigningMethodHS256, security.Claims{
			UserID: 1,
			Email:  "user@email.com",
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			},
		})

		claims, err := security.ValidateToken(token)

		assert.Nil(t, claims)
		assert.Error(t, err)
	})

	t.Run("should return error when token algorithm is invalid", func(t *testing.T) {
		setupJWTEnv(t)

		token := makeToken(t, jwt.SigningMethodHS384, security.Claims{
			UserID: 1,
			Email:  "user@email.com",
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		claims, err := security.ValidateToken(token)

		assert.Nil(t, claims)
		assert.Error(t, err)
	})

	t.Run("should return error when token claims are invalid", func(t *testing.T) {
		setupJWTEnv(t)

		token := makeToken(t, jwt.SigningMethodHS256, security.Claims{
			UserID: 0,
			Email:  " ",
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		})

		claims, err := security.ValidateToken(token)

		assert.Nil(t, claims)
		assert.EqualError(t, err, "claims do token inv\u00e1lidas")
	})
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		expectedToken string
		expectedErr   string
	}{
		{
			name:          "should extract token from bearer header",
			authorization: "Bearer abc.def.ghi",
			expectedToken: "abc.def.ghi",
		},
		{
			name:          "should trim spaces and ignore bearer case",
			authorization: "  bearer   abc.def.ghi  ",
			expectedToken: "abc.def.ghi",
		},
		{
			name:          "should return error when header is empty",
			authorization: " ",
			expectedErr:   "header Authorization n\u00e3o informado",
		},
		{
			name:          "should return error when header format is invalid",
			authorization: "Token abc.def.ghi",
			expectedErr:   "header Authorization deve estar no formato Bearer token",
		},
		{
			name:          "should return error when token is empty",
			authorization: "Bearer ",
			expectedErr:   "header Authorization deve estar no formato Bearer token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := security.ExtractBearerToken(tt.authorization)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}

func TestGetExpirationDuration(t *testing.T) {
	t.Run("should return configured expiration duration", func(t *testing.T) {
		setupJWTEnv(t)
		env.AuthTokenExpirationTimeInHours = 1.5

		duration := security.GetExpirationDuration()

		assert.Equal(t, 90*time.Minute, duration)
	})

	t.Run("should return default expiration duration when configured value is invalid", func(t *testing.T) {
		setupJWTEnv(t)
		env.AuthTokenExpirationTimeInHours = 0

		duration := security.GetExpirationDuration()

		assert.Equal(t, 8*time.Hour, duration)
	})
}

func setupJWTEnv(t *testing.T) {
	t.Helper()

	originalJWTSecret := env.JWTSecret
	originalAppName := env.AppName
	originalExpiration := env.AuthTokenExpirationTimeInHours

	env.JWTSecret = "test-secret"
	env.AppName = "streaming_api_test"
	env.AuthTokenExpirationTimeInHours = 2

	t.Cleanup(func() {
		env.JWTSecret = originalJWTSecret
		env.AppName = originalAppName
		env.AuthTokenExpirationTimeInHours = originalExpiration
	})
}

func makeToken(t *testing.T, method jwt.SigningMethod, claims security.Claims) string {
	t.Helper()

	token, err := jwt.NewWithClaims(method, claims).SignedString([]byte(env.JWTSecret))

	require.NoError(t, err)
	return token
}
