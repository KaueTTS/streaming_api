package env_test

import (
	"testing"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setDefaultEnv(t *testing.T) {
	t.Helper()

	t.Setenv("PORT", "8080")
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_NAME", "test_app")
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
	t.Setenv("SQLITE_DATABASE_URL", "file:test.db?cache=shared&mode=memory")
	t.Setenv("JWT_SECRET", "test_secret")
	t.Setenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS", "4")
	t.Setenv("REDIS_HOST", "redis")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_USERNAME", "default")
	t.Setenv("REDIS_PASSWORD", "redis_password")
	t.Setenv("REDIS_DATABASE", "1")
	t.Setenv("TMDB_BASE_URL", "https://api.themoviedb.org/3")
	t.Setenv("TMDB_ACCESS_TOKEN", "test_tmdb_access_token")
}

func TestInit(t *testing.T) {
	t.Run("success in using environment values", func(t *testing.T) {
		setDefaultEnv(t)

		err := env.Init()

		require.NoError(t, err)
		assert.Equal(t, "8080", env.Port)
		assert.Equal(t, "test", env.AppEnv)
		assert.Equal(t, "test_app", env.AppName)
		assert.Equal(t, "localhost:4318", env.OTLPExporterEndpoint)
		assert.Equal(t, "file:test.db?cache=shared&mode=memory", env.SQLiteDatabaseURL)
		assert.Equal(t, "test_secret", env.JWTSecret)
		assert.Equal(t, 4.0, env.AuthTokenExpirationTimeInHours)
		assert.Equal(t, "redis", env.RedisHost)
		assert.Equal(t, "6379", env.RedisPort)
		assert.Equal(t, "default", env.RedisUsername)
		assert.Equal(t, "redis_password", env.RedisPassword)
		assert.Equal(t, 1, env.RedisDatabase)
		assert.Equal(t, "https://api.themoviedb.org/3", env.TMDBBaseURL)
		assert.Equal(t, "test_tmdb_access_token", env.TMDBAccessToken)
	})

	tests := []struct {
		name        string
		envKey      string
		envValue    string
		expectedErr string
	}{
		{
			name:        "should return error if JWT_SECRET is missing",
			envKey:      "JWT_SECRET",
			envValue:    "",
			expectedErr: "a variável de ambiente JWT_SECRET precisa ser informada",
		},
		{
			name:        "should return error if SQLITE_DATABASE_URL is missing",
			envKey:      "SQLITE_DATABASE_URL",
			envValue:    "",
			expectedErr: "a variável de ambiente SQLITE_DATABASE_URL precisa ser informada",
		},
		{
			name:        "should return error if redis database is invalid",
			envKey:      "REDIS_DATABASE",
			envValue:    "-1",
			expectedErr: "REDIS_DATABASE precisa ser um número maior ou igual a zero",
		},
		{
			name:        "should return error if auth token expiration is invalid",
			envKey:      "AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS",
			envValue:    "0",
			expectedErr: "AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS precisa ser um número maior que zero",
		},
		{
			name:        "should return error if TMDB_BASE_URL is missing",
			envKey:      "TMDB_BASE_URL",
			envValue:    "",
			expectedErr: "a variável de ambiente TMDB_BASE_URL precisa ser informada",
		},
		{
			name:        "should return error if TMDB_ACCESS_TOKEN is missing",
			envKey:      "TMDB_ACCESS_TOKEN",
			envValue:    "",
			expectedErr: "a variável de ambiente TMDB_ACCESS_TOKEN precisa ser informada",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setDefaultEnv(t)
			t.Setenv(tt.envKey, tt.envValue)

			err := env.Init()

			require.Error(t, err)
			assert.EqualError(t, err, tt.expectedErr)
		})
	}
}
