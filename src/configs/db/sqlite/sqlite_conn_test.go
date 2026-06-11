package sqlite_conn_test

import (
	"os"
	"path/filepath"
	"testing"

	sqlite_conn "github.com/KaueTTS/streaming_api/src/configs/db/sqlite"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	models "github.com/KaueTTS/streaming_api/src/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("should be successfully initialized", func(t *testing.T) {
		tempDir := t.TempDir()

		oldDatabaseURL := env.SQLiteDatabaseURL
		env.SQLiteDatabaseURL = filepath.Join(tempDir, "test.db")

		t.Cleanup(func() {
			env.SQLiteDatabaseURL = oldDatabaseURL
		})

		db, err := sqlite_conn.Init()
		require.NoError(t, err)
		require.NotNil(t, db)

		sqlDB, err := db.DB()
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, sqlDB.Close())
		})

		assert.FileExists(t, env.SQLiteDatabaseURL)
		assert.True(t, db.Migrator().HasTable(&models.User{}))
		assert.True(t, db.Migrator().HasTable(&models.Profile{}))
		assert.True(t, db.Migrator().HasTable(&models.Content{}))
		assert.True(t, db.Migrator().HasTable(&models.Favorite{}))
		assert.True(t, db.Migrator().HasTable(&models.WatchProgress{}))
	})

	t.Run("should return error when database directory cannot be created", func(t *testing.T) {
		tempDir := t.TempDir()
		blockingFilePath := filepath.Join(tempDir, "database")

		require.NoError(t, os.WriteFile(blockingFilePath, []byte("not a directory"), 0o600))

		oldDatabaseURL := env.SQLiteDatabaseURL
		env.SQLiteDatabaseURL = filepath.Join(blockingFilePath, "test.db")

		t.Cleanup(func() {
			env.SQLiteDatabaseURL = oldDatabaseURL
		})

		db, err := sqlite_conn.Init()

		require.Error(t, err)
		assert.Nil(t, db)
		assert.ErrorContains(t, err, "erro ao criar diretório do banco de dados")
	})
}
