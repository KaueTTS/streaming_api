package redis_conn

import (
	"context"
	"fmt"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	"github.com/redis/go-redis/v9"
)

func Init(ctx context.Context) (*redis.Client, error) {
	if env.RedisHost == "" {
		return nil, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.RedisHost, env.RedisPort),
		Username: env.RedisUsername,
		Password: env.RedisPassword,
		DB:       env.RedisDatabase,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("erro ao conectar no redis: %w", err)
	}

	return client, nil
}
