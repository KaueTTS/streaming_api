package main

import (
	"fmt"
	"log"

	_ "github.com/KaueTTS/streaming_api/docs"
	api "github.com/KaueTTS/streaming_api/src/api"
	sqlite_conn "github.com/KaueTTS/streaming_api/src/configs/db/sqlite"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
)

// @title Streaming API
// @version 1.0
// @description API

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
	if err := env.Init(); err != nil {
		return fmt.Errorf("erro ao inicializar variáveis de ambiente: %w", err)
	}

	db, err := sqlite_conn.Init()
	if err != nil {
		return fmt.Errorf("erro ao inicializar sqlite: %w", err)
	}

	if err := api.Init(db); err != nil {
		return fmt.Errorf("erro ao iniciar api: %w", err)
	}

	return nil
}
