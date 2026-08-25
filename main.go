package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	// Configurar o formato do log para JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Cria um contexto que é cancelado automaticamente se o app receber um sinal de parar
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Executa a aplicação
	if err := run(ctx); err != nil {
		slog.Error("Falha ao iniciar aplicação", "erro", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Inicialização das variáveis
	if err := env.Init(); err != nil {
		return fmt.Errorf("erro ao inicializar variáveis de ambiente: %w", err)
	}

	// Inicialização do tracing
	tracerProvider, err := tracing.Init(ctx)
	if err != nil {
		return fmt.Errorf("erro ao inicializar tracing: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			slog.Error("Erro ao finalizar tracing", "erro", err)
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
		slog.Warn("Redis indisponível, usando fallback em memória", "erro", err)
		redisClient = nil
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	// Inicialização da API
	app := api.Init(db, redisClient)
	if err := api.Listen(ctx, app); err != nil {
		return fmt.Errorf("erro ao iniciar api: %w", err)
	}

	return nil
}
