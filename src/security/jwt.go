package security

import (
	"errors"
	"strings"
	"time"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, role string) (string, error) {
	secret := shared_normalizers.TrimString(env.JWTSecret)
	if secret == "" {
		return "", errors.New("JWT_SECRET não informado")
	}

	now := time.Now()
	expiresAt := now.Add(GetExpirationDuration())

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.AppName,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string) (*Claims, error) {
	secret := shared_normalizers.TrimString(env.JWTSecret)
	if secret == "" {
		return nil, errors.New("JWT_SECRET não informado")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("algoritmo do token inválido")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	if claims.UserID == 0 || shared_normalizers.TrimString(claims.Email) == "" {
		return nil, errors.New("claims do token inválidas")
	}

	return claims, nil
}

func ExtractBearerToken(authorizationHeader string) (string, error) {
	authorizationHeader = shared_normalizers.TrimString(authorizationHeader)
	if authorizationHeader == "" {
		return "", errors.New("header Authorization não informado")
	}

	parts := strings.SplitN(authorizationHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("header Authorization deve estar no formato Bearer token")
	}

	token := shared_normalizers.TrimString(parts[1])
	if token == "" {
		return "", errors.New("token não informado")
	}

	return token, nil
}

func GetExpirationDuration() time.Duration {
	if env.AuthTokenExpirationTimeInHours <= 0 {
		return 8 * time.Hour
	}

	return time.Duration(env.AuthTokenExpirationTimeInHours * float64(time.Hour))
}
