package repository_sqlite_user_test

import (
	"context"
	"testing"

	models "github.com/KaueTTS/streaming_api/src/models"
	repository_sqlite_user "github.com/KaueTTS/streaming_api/src/repositories/sqlite/user"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	t.Run("should create user", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)
		user := makeUser("John Doe", "john@example.com")

		err := repository.Create(context.Background(), user)

		require.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.False(t, user.CreatedAt.IsZero())
		assert.False(t, user.UpdatedAt.IsZero())

		var savedUser models.User
		err = db.First(&savedUser, user.ID).Error

		require.NoError(t, err)
		assert.Equal(t, user.Name, savedUser.Name)
		assert.Equal(t, user.Email, savedUser.Email)
		assert.Equal(t, user.PasswordHash, savedUser.PasswordHash)
		assert.Equal(t, user.Role, savedUser.Role)
	})

	t.Run("should return error when email already exists", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)
		firstUser := makeUser("John Doe", "john@example.com")
		secondUser := makeUser("Jane Doe", "john@example.com")

		err := repository.Create(context.Background(), firstUser)
		require.NoError(t, err)

		err = repository.Create(context.Background(), secondUser)

		assert.Error(t, err)
	})
}

func TestFindByEmail(t *testing.T) {
	t.Run("should return user by normalized email", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)
		user := makeUser("John Doe", "john@example.com")

		err := repository.Create(context.Background(), user)
		require.NoError(t, err)

		foundUser, err := repository.FindByEmail(context.Background(), "  JOHN@EXAMPLE.COM  ")

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.PasswordHash, foundUser.PasswordHash)
		assert.Equal(t, user.Role, foundUser.Role)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)

		foundUser, err := repository.FindByEmail(context.Background(), "missing@example.com")

		assert.Nil(t, foundUser)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestFindByID(t *testing.T) {
	t.Run("should return user by id", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)
		user := makeUser("John Doe", "john@example.com")

		err := repository.Create(context.Background(), user)
		require.NoError(t, err)

		foundUser, err := repository.FindByID(context.Background(), user.ID)

		require.NoError(t, err)
		require.NotNil(t, foundUser)
		assert.Equal(t, user.ID, foundUser.ID)
		assert.Equal(t, user.Name, foundUser.Name)
		assert.Equal(t, user.Email, foundUser.Email)
		assert.Equal(t, user.PasswordHash, foundUser.PasswordHash)
		assert.Equal(t, user.Role, foundUser.Role)
	})

	t.Run("should return error when user is not found", func(t *testing.T) {
		db := setupUserRepositoryTestDB(t)
		repository := repository_sqlite_user.NewUserRepository(db)

		foundUser, err := repository.FindByID(context.Background(), 999)

		assert.Nil(t, foundUser)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func setupUserRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func makeUser(name string, email string) *models.User {
	return &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: "hashed-password",
		Role:         "user",
	}
}
