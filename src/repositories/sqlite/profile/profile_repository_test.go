package repository_sqlite_profile_test

import (
	"context"
	"testing"

	models "github.com/KaueTTS/streaming_api/src/models"
	repository_sqlite_profile "github.com/KaueTTS/streaming_api/src/repositories/sqlite/profile"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestFindByUserID(t *testing.T) {
	t.Run("should list only profiles from selected user", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createProfileTestUser(t, db, "Jane Doe", "jane@example.com")
		firstProfile := createProfileTestProfile(t, db, user.ID, "Main", false)
		secondProfile := createProfileTestProfile(t, db, user.ID, "Kids", true)
		createProfileTestProfile(t, db, otherUser.ID, "Other User Profile", false)

		profiles, total, err := repository.FindByUserID(context.Background(), user.ID, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		require.Len(t, profiles, 2)

		profilesByName := map[string]models.Profile{}
		for _, profile := range profiles {
			profilesByName[profile.Name] = profile
			assert.Equal(t, user.ID, profile.UserID)
		}

		assert.Equal(t, firstProfile.ID, profilesByName["Main"].ID)
		assert.Equal(t, false, profilesByName["Main"].IsKids)
		assert.Equal(t, secondProfile.ID, profilesByName["Kids"].ID)
		assert.Equal(t, true, profilesByName["Kids"].IsKids)
		assert.NotContains(t, profilesByName, "Other User Profile")
	})

	t.Run("should return empty list when user has no profiles", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")

		profiles, total, err := repository.FindByUserID(context.Background(), user.ID, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.Empty(t, profiles)
	})

	t.Run("should paginate results", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		createProfileTestProfile(t, db, user.ID, "Alpha", false)
		createProfileTestProfile(t, db, user.ID, "Beta", false)
		createProfileTestProfile(t, db, user.ID, "Gamma", false)

		profiles, total, err := repository.FindByUserID(context.Background(), user.ID, 1, 2)

		require.NoError(t, err)
		assert.Equal(t, int64(3), total)
		assert.Len(t, profiles, 2)
		// ordered by name asc: Alpha, Beta
		assert.Equal(t, "Alpha", profiles[0].Name)
		assert.Equal(t, "Beta", profiles[1].Name)
	})
}

func TestCountByUserID(t *testing.T) {
	t.Run("should return correct count of profiles for user", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createProfileTestUser(t, db, "Jane Doe", "jane@example.com")
		createProfileTestProfile(t, db, user.ID, "Main", false)
		createProfileTestProfile(t, db, user.ID, "Kids", true)
		createProfileTestProfile(t, db, otherUser.ID, "Other", false)

		count, err := repository.CountByUserID(context.Background(), user.ID)

		require.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("should return zero when user has no profiles", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")

		count, err := repository.CountByUserID(context.Background(), user.ID)

		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestCreate(t *testing.T) {
	t.Run("should create profile", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		profile := &models.Profile{
			UserID: user.ID,
			Name:   "Main",
			IsKids: false,
		}

		err := repository.Create(context.Background(), profile)

		require.NoError(t, err)
		assert.NotZero(t, profile.ID)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())

		var savedProfile models.Profile
		err = db.First(&savedProfile, profile.ID).Error

		require.NoError(t, err)
		assert.Equal(t, profile.UserID, savedProfile.UserID)
		assert.Equal(t, profile.Name, savedProfile.Name)
		assert.Equal(t, profile.IsKids, savedProfile.IsKids)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("should update profile fields", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		profile := createProfileTestProfile(t, db, user.ID, "Main", false)

		avatarURL := "https://example.com/avatar.png"
		profile.Name = "Updated Name"
		profile.AvatarURL = &avatarURL
		profile.IsKids = true

		err := repository.Update(context.Background(), &profile)

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", profile.Name)
		assert.Equal(t, &avatarURL, profile.AvatarURL)
		assert.Equal(t, true, profile.IsKids)

		var savedProfile models.Profile
		err = db.First(&savedProfile, profile.ID).Error

		require.NoError(t, err)
		assert.Equal(t, "Updated Name", savedProfile.Name)
		assert.Equal(t, true, savedProfile.IsKids)
	})

	t.Run("should return error when profile does not belong to user", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createProfileTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createProfileTestProfile(t, db, otherUser.ID, "Main", false)

		profile.UserID = user.ID
		profile.Name = "Hacked"

		err := repository.Update(context.Background(), &profile)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestDelete(t *testing.T) {
	t.Run("should delete profile", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		profile := createProfileTestProfile(t, db, user.ID, "Main", false)

		err := repository.Delete(context.Background(), user.ID, profile.ID)

		require.NoError(t, err)

		var savedProfile models.Profile
		err = db.Unscoped().First(&savedProfile, profile.ID).Error
		require.NoError(t, err)
		assert.True(t, savedProfile.DeletedAt.Valid)
	})

	t.Run("should return error when profile does not belong to user", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		otherUser := createProfileTestUser(t, db, "Jane Doe", "jane@example.com")
		profile := createProfileTestProfile(t, db, otherUser.ID, "Main", false)

		err := repository.Delete(context.Background(), user.ID, profile.ID)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)
		user := createProfileTestUser(t, db, "John Doe", "john@example.com")

		err := repository.Delete(context.Background(), user.ID, 999)

		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestFindProfileByID(t *testing.T) {
	t.Run("should list only selected profile", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)

		user := createProfileTestUser(t, db, "John Doe", "john@example.com")
		profile := createProfileTestProfile(t, db, user.ID, "Main", false)

		foundProfile, err := repository.FindProfileByID(
			context.Background(),
			profile.ID,
		)

		require.NoError(t, err)
		require.NotNil(t, foundProfile)

		assert.Equal(t, profile.ID, foundProfile.ID)
		assert.Equal(t, profile.UserID, foundProfile.UserID)
		assert.Equal(t, profile.Name, foundProfile.Name)
		assert.Equal(t, profile.IsKids, foundProfile.IsKids)
	})

	t.Run("should return error when profile does not exist", func(t *testing.T) {
		db := setupProfileRepositoryTestDB(t)
		repository := repository_sqlite_profile.NewProfileRepository(db)

		foundProfile, err := repository.FindProfileByID(
			context.Background(),
			999,
		)

		assert.Nil(t, foundProfile)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func setupProfileRepositoryTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		TranslateError: true,
	})
	require.NoError(t, err)
	require.NotNil(t, db)

	err = db.AutoMigrate(&models.User{}, &models.Profile{})
	require.NoError(t, err)

	return db
}

func createProfileTestUser(t *testing.T, db *gorm.DB, name string, email string) models.User {
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

func createProfileTestProfile(t *testing.T, db *gorm.DB, userID uint, name string, isKids bool) models.Profile {
	t.Helper()

	profile := models.Profile{
		UserID: userID,
		Name:   name,
		IsKids: isKids,
	}
	require.NoError(t, db.Create(&profile).Error)

	return profile
}
