package repository_sqlite_favorite_test

import (
	"context"
	"testing"
	"time"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	models "github.com/KaueTTS/streaming_api/src/models"
	repository_sqlite_favorite "github.com/KaueTTS/streaming_api/src/repositories/sqlite/favorite"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFindFavoriteByProfileID(t *testing.T) {
	t.Run("should list only favorites from selected profile", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createFavoriteTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		otherProfile := createFavoriteTestProfile(t, db, otherUser.ID, "Other User Profile", false)
		firstFavorite := createFavoriteTestFavorite(t, db, user.ID, profile.ID, 101, "movie", time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC))
		secondFavorite := createFavoriteTestFavorite(t, db, user.ID, profile.ID, 202, "tv", time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC))
		createFavoriteTestFavorite(t, db, otherUser.ID, otherProfile.ID, 303, "movie", time.Date(2026, 1, 3, 10, 0, 0, 0, time.UTC))

		favorites, total, err := repository.FindFavoriteByProfileID(context.Background(), user.ID, profile.ID, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, favorites, 2)
		assert.Equal(t, secondFavorite.ID, favorites[0].ID)
		assert.Equal(t, secondFavorite.ContentExternalId, favorites[0].ContentExternalId)
		assert.Equal(t, firstFavorite.ID, favorites[1].ID)
		assert.Equal(t, firstFavorite.ContentExternalId, favorites[1].ContentExternalId)

		for _, favorite := range favorites {
			assert.Equal(t, profile.ID, favorite.ProfileID)
		}
	})

	t.Run("should paginate results", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		createFavoriteTestFavorite(t, db, user.ID, profile.ID, 101, "movie", time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC))
		expectedFavorite := createFavoriteTestFavorite(t, db, user.ID, profile.ID, 202, "tv", time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC))
		createFavoriteTestFavorite(t, db, user.ID, profile.ID, 303, "movie", time.Date(2026, 1, 3, 10, 0, 0, 0, time.UTC))

		favorites, total, err := repository.FindFavoriteByProfileID(context.Background(), user.ID, profile.ID, 2, 1)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		require.Len(t, favorites, 1)
		assert.Equal(t, expectedFavorite.ID, favorites[0].ID)
		assert.Equal(t, expectedFavorite.ContentExternalId, favorites[0].ContentExternalId)
	})

	t.Run("should return error when profile does not belong to user", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createFavoriteTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createFavoriteTestProfile(t, db, otherUser.ID, "Other User Profile", false)

		favorites, total, err := repository.FindFavoriteByProfileID(context.Background(), user.ID, profile.ID, 1, 10)

		assert.Nil(t, favorites)
		assert.Equal(t, int64(0), total)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestCreateFavoriteByProfileID(t *testing.T) {
	t.Run("should create favorite", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: 101,
			Type:              "movie",
		}

		err := repository.CreateFavoriteByProfileID(context.Background(), user.ID, request)

		require.NoError(t, err)

		var savedFavorite models.Favorite
		err = db.First(&savedFavorite).Error

		require.NoError(t, err)
		assert.NotZero(t, savedFavorite.ID)
		assert.Equal(t, user.ID, savedFavorite.UserID)
		assert.Equal(t, profile.ID, savedFavorite.ProfileID)
		assert.Equal(t, request.ContentExternalID, savedFavorite.ContentExternalId)
		assert.Equal(t, request.Type, savedFavorite.Type)
		assert.False(t, savedFavorite.CreatedAt.IsZero())
		assert.False(t, savedFavorite.UpdatedAt.IsZero())
	})

	t.Run("should return error when profile does not belong to user", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createFavoriteTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createFavoriteTestProfile(t, db, otherUser.ID, "Other User Profile", false)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: 101,
			Type:              "movie",
		}

		err := repository.CreateFavoriteByProfileID(context.Background(), user.ID, request)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("should return error when favorite already exists", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: 101,
			Type:              "movie",
		}

		err := repository.CreateFavoriteByProfileID(context.Background(), user.ID, request)
		require.NoError(t, err)

		err = repository.CreateFavoriteByProfileID(context.Background(), user.ID, request)

		assert.Error(t, err)
	})
}

func TestDeleteFavoriteByProfileID(t *testing.T) {
	t.Run("should delete favorite", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		favorite := createFavoriteTestFavorite(t, db, user.ID, profile.ID, 101, "movie", time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC))
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: favorite.ContentExternalId,
			Type:              favorite.Type,
		}

		err := repository.DeleteFavoriteByProfileID(context.Background(), user.ID, request)

		require.NoError(t, err)

		var savedFavorite models.Favorite
		err = db.Unscoped().First(&savedFavorite, favorite.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("should return error when profile does not belong to user", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createFavoriteTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createFavoriteTestProfile(t, db, otherUser.ID, "Other User Profile", false)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: 101,
			Type:              "movie",
		}

		err := repository.DeleteFavoriteByProfileID(context.Background(), user.ID, request)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("should return error when favorite does not exist", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main", false)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         profile.ID,
			ContentExternalID: 999,
			Type:              "movie",
		}

		err := repository.DeleteFavoriteByProfileID(context.Background(), user.ID, request)

		assert.ErrorIs(t, err, shared_errors.ErrFavoriteNotFound)
	})
}

func setupFavoriteRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.AutoMigrate(&models.User{}, &models.Profile{}, &models.Favorite{})
	require.NoError(t, err)

	return db
}

func createFavoriteTestUser(t *testing.T, db *gorm.DB, name string, email string) models.User {
	t.Helper()

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: "hashed-password",
		Role:         "user",
	}
	require.NoError(t, db.Create(&user).Error)

	return user
}

func createFavoriteTestProfile(t *testing.T, db *gorm.DB, userID uint, name string, isKids bool) models.Profile {
	t.Helper()

	profile := models.Profile{
		UserID: userID,
		Name:   name,
		IsKids: isKids,
	}
	require.NoError(t, db.Create(&profile).Error)

	return profile
}

func createFavoriteTestFavorite(t *testing.T, db *gorm.DB, userID uint, profileID uint, contentExternalID int, contentType string, createdAt time.Time) models.Favorite {
	t.Helper()

	favorite := models.Favorite{
		UserID:            userID,
		ProfileID:         profileID,
		ContentExternalId: contentExternalID,
		Type:              contentType,
		CreatedAt:         createdAt,
	}
	require.NoError(t, db.Create(&favorite).Error)

	return favorite
}
