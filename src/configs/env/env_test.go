package env_test

import (
	"testing"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("success in using environmental values", func(t *testing.T) {
		t.Setenv("PORT", "8080")
		t.Setenv("APP_ENV", "test")
		t.Setenv("APP_NAME", "test_app")
		t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318")
		t.Setenv("SQLITE_DATABASE_URL", "file:test.db?cache=shared&mode=memory")
		t.Setenv("JWT_SECRET", "test_secret")
		t.Setenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS", "4")

		err := env.Init()

		assert.NoError(t, err)
		assert.Equal(t, "8080", env.Port)
		assert.Equal(t, "test", env.AppEnv)
		assert.Equal(t, "test_app", env.AppName)
		assert.Equal(t, "localhost:4318", env.OTLPExporterEndpoint)
		assert.Equal(t, "file:test.db?cache=shared&mode=memory", env.SQLiteDatabaseURL)
		assert.Equal(t, "test_secret", env.JWTSecret)
		assert.Equal(t, 4.0, env.AuthTokenExpirationTimeInHours)
	})

	t.Run("should return error if SQLITE_DATABASE_URL is missing", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "test_secret")
		t.Setenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS", "4")
		t.Setenv("SQLITE_DATABASE_URL", "")

		err := env.Init()

		assert.Error(t, err)
		assert.Equal(t, "a variável de ambiente SQLITE_DATABASE_URL precisa ser informada", err.Error())

	})

	t.Run("should return error if JWT_SECRET is missing", func(t *testing.T) {
		t.Setenv("SQLITE_DATABASE_URL", "file:test.db?cache=shared&mode=memory")
		t.Setenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS", "4")
		t.Setenv("JWT_SECRET", "")

		err := env.Init()

		assert.Error(t, err)
		assert.Equal(t, "a variável de ambiente JWT_SECRET precisa ser informada", err.Error())
	})

	t.Run("should return error if auth token expiration is invalid", func(t *testing.T) {
		t.Setenv("SQLITE_DATABASE_URL", "file:test.db?cache=shared&mode=memory")
		t.Setenv("JWT_SECRET", "test_secret")
		t.Setenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS", "0")

		err := env.Init()

		assert.Error(t, err)
		assert.Equal(t, "AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS precisa ser um número maior que zero", err.Error())
	})
}
