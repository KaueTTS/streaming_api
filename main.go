package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/KaueTTS/streaming_api/docs"
	api "github.com/KaueTTS/streaming_api/src/api"
	redis_conn "github.com/KaueTTS/streaming_api/src/configs/db/redis"
	sqlite_conn "github.com/KaueTTS/streaming_api/src/configs/db/sqlite"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	tracing "github.com/KaueTTS/streaming_api/src/configs/tracing"
)

// @title Streaming API
// @version 1.0
// @description API de streaming

// @contact.name KauêTTS
// @contact.url https://github.com/KaueTTS

// @accept json
// @produce json

// @schemes http https

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := run(); err != nil {
		log.Fatalf("falha ao iniciar aplicação: %v", err)
	}
}

func run() error {
	// Inicialização das variáveis
	if err := env.Init(); err != nil {
		return fmt.Errorf("erro ao inicializar variáveis de ambiente: %w", err)
	}

	// Inicialização do tracing
	ctx := context.Background()
	tracerProvider, err := tracing.Init(ctx)
	if err != nil {
		return fmt.Errorf("erro ao inicializar tracing: %w", err)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("erro ao finalizar tracing: %v", err)
		}
	}()

	// Inicialização do SQLite
	db, err := sqlite_conn.Init()
	if err != nil {
		return fmt.Errorf("erro ao inicializar sqlite: %w", err)
	}

	// Inicialização do Redis
	redisClient, err := redis_conn.Init(ctx)
	if err != nil {
		log.Printf("Redis indisponível, usando fallback em memória: %v", err)
		redisClient = nil
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	// Inicialização da API
	app := api.Init(db, redisClient)
	if err := api.Listen(app); err != nil {
		return fmt.Errorf("erro ao iniciar api: %w", err)
	}

	return nil
}
