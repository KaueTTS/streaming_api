package redis_conn_test

import (
	"context"
	"testing"

	redis_conn "github.com/KaueTTS/streaming_api/src/configs/db/redis"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("should be successfully initialized", func(t *testing.T) {
		redisServer := miniredis.RunT(t)

		oldRedisHost := env.RedisHost
		oldRedisPort := env.RedisPort
		oldRedisUsername := env.RedisUsername
		oldRedisPassword := env.RedisPassword
		oldRedisDatabase := env.RedisDatabase

		env.RedisHost = redisServer.Host()
		env.RedisPort = redisServer.Port()
		env.RedisUsername = ""
		env.RedisPassword = ""
		env.RedisDatabase = 0

		t.Cleanup(func() {
			env.RedisHost = oldRedisHost
			env.RedisPort = oldRedisPort
			env.RedisUsername = oldRedisUsername
			env.RedisPassword = oldRedisPassword
			env.RedisDatabase = oldRedisDatabase
		})

		client, err := redis_conn.Init(context.Background())

		require.NoError(t, err)
		require.NotNil(t, client)

		t.Cleanup(func() {
			require.NoError(t, client.Close())
		})

		assert.NoError(t, client.Ping(context.Background()).Err())
	})

	t.Run("should return nil when redis host is empty", func(t *testing.T) {
		oldRedisHost := env.RedisHost
		oldRedisPort := env.RedisPort
		oldRedisUsername := env.RedisUsername
		oldRedisPassword := env.RedisPassword
		oldRedisDatabase := env.RedisDatabase

		env.RedisHost = ""
		env.RedisPort = ""
		env.RedisUsername = ""
		env.RedisPassword = ""
		env.RedisDatabase = 0

		t.Cleanup(func() {
			env.RedisHost = oldRedisHost
			env.RedisPort = oldRedisPort
			env.RedisUsername = oldRedisUsername
			env.RedisPassword = oldRedisPassword
			env.RedisDatabase = oldRedisDatabase
		})

		client, err := redis_conn.Init(context.Background())

		require.NoError(t, err)
		assert.Nil(t, client)
	})

	t.Run("should return error when redis is unavailable", func(t *testing.T) {
		oldRedisHost := env.RedisHost
		oldRedisPort := env.RedisPort
		oldRedisUsername := env.RedisUsername
		oldRedisPassword := env.RedisPassword
		oldRedisDatabase := env.RedisDatabase

		env.RedisHost = "127.0.0.1"
		env.RedisPort = "0"
		env.RedisUsername = ""
		env.RedisPassword = ""
		env.RedisDatabase = 0

		t.Cleanup(func() {
			env.RedisHost = oldRedisHost
			env.RedisPort = oldRedisPort
			env.RedisUsername = oldRedisUsername
			env.RedisPassword = oldRedisPassword
			env.RedisDatabase = oldRedisDatabase
		})

		client, err := redis_conn.Init(context.Background())

		require.Error(t, err)
		assert.Nil(t, client)
		assert.ErrorContains(t, err, "erro ao conectar no redis")
	})
}
