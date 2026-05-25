package sqlite_conn

import (
	"fmt"
	"os"
	"path/filepath"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	models "github.com/KaueTTS/streaming_api/src/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dir := filepath.Dir(env.SQLiteDatabaseURL)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório do banco de dados: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(env.SQLiteDatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no sqlite: %w", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Content{},
		&models.Favorite{},
		&models.WatchProgress{},
	); err != nil {
		return nil, fmt.Errorf("erro ao migrar banco de dados: %w", err)
	}

	return db, nil
}
