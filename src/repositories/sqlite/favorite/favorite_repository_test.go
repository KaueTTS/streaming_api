package repository_sqlite_favorite_test

import (
	"context"
	"testing"

	models "github.com/KaueTTS/streaming_api/src/models"
	repository_sqlite_favorite "github.com/KaueTTS/streaming_api/src/repositories/sqlite/favorite"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFindByProfileID(t *testing.T) {
	t.Run("should list only favorites from selected profile", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		profile := createFavoriteTestProfile(t, db, user.ID, "Main")
		otherProfile := createFavoriteTestProfile(t, db, user.ID, "Kids")
		firstFavorite := createFavoriteTestFavorite(t, db, profile.ID, 100, "movie")
		secondFavorite := createFavoriteTestFavorite(t, db, profile.ID, 200, "tv")
		createFavoriteTestFavorite(t, db, otherProfile.ID, 300, "movie")

		favorites, total, err := repository.FindByProfileID(context.Background(), user.ID, profile.ID, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, favorites, 2)

		favoritesByContentID := map[int]models.Favorite{}
		for _, favorite := range favorites {
			favoritesByContentID[favorite.ContentID] = favorite
			assert.Equal(t, profile.ID, favorite.ProfileID)
		}

		assert.Equal(t, firstFavorite.ID, favoritesByContentID[100].ID)
		assert.Equal(t, "movie", favoritesByContentID[100].ContentType)
		assert.Equal(t, secondFavorite.ID, favoritesByContentID[200].ID)
		assert.Equal(t, "tv", favoritesByContentID[200].ContentType)
		assert.NotContains(t, favoritesByContentID, 300)
	})

	t.Run("should not list favorites from profile owned by another user", func(t *testing.T) {
		db := setupFavoriteRepositoryTestDB(t)
		repository := repository_sqlite_favorite.NewFavoriteRepository(db)
		user := createFavoriteTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createFavoriteTestUser(t, db, "Jane Doe", "jane@example.com")
		otherProfile := createFavoriteTestProfile(t, db, otherUser.ID, "Main")
		createFavoriteTestFavorite(t, db, otherProfile.ID, 100, "movie")

		favorites, total, err := repository.FindByProfileID(context.Background(), user.ID, otherProfile.ID, 1, 10)

		assert.Nil(t, favorites)
		assert.Zero(t, total)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func setupFavoriteRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

func createFavoriteTestProfile(t *testing.T, db *gorm.DB, userID uint, name string) models.Profile {
	t.Helper()

	profile := models.Profile{
		UserID: userID,
		Name:   name,
		IsKids: false,
	}
	require.NoError(t, db.Create(&profile).Error)

	return profile
}

func createFavoriteTestFavorite(t *testing.T, db *gorm.DB, profileID uint, contentID int, contentType string) models.Favorite {
	t.Helper()

	favorite := models.Favorite{
		ProfileID:   profileID,
		ContentID:   contentID,
		ContentType: contentType,
	}
	require.NoError(t, db.Create(&favorite).Error)

	return favorite
}
