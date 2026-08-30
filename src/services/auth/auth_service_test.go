package service_auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	repository_mock "github.com/KaueTTS/streaming_api/src/mocks/repositories"
	models "github.com/KaueTTS/streaming_api/src/models"
	security "github.com/KaueTTS/streaming_api/src/security"
	service_auth "github.com/KaueTTS/streaming_api/src/services/auth"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {
	t.Run("should register user with normalized data", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		request := dto_auth.RegisterRequestDto{
			Name:     "  John Doe  ",
			Email:    "  JOHN@EXAMPLE.COM  ",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return((*models.User)(nil), gorm.ErrRecordNotFound).
			Once()
		repository.On("Create", ctx, mock.MatchedBy(func(user *models.User) bool {
			return user.Name == "John Doe" &&
				user.Email == "john@example.com" &&
				user.Role == shared_constants.RoleUser &&
				user.PasswordHash != "" &&
				user.PasswordHash != request.Password &&
				security.Compare(user.PasswordHash, request.Password)
		})).
			Run(func(args mock.Arguments) {
				user := args.Get(1).(*models.User)
				user.ID = 1
				user.CreatedAt = createdAt
				user.UpdatedAt = updatedAt
			}).
			Return(nil).
			Once()

		response, err := service.Register(ctx, request)

		require.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "John Doe", response.Name)
		assert.Equal(t, "john@example.com", response.Email)
		assert.Equal(t, shared_constants.RoleUser, response.Role)
		assert.Equal(t, createdAt, response.CreatedAt)
		assert.Equal(t, updatedAt, response.UpdatedAt)
		repository.AssertExpectations(t)
	})

	t.Run("should return invalid password when password validation fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		request := dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "short1",
		}

		response, err := service.Register(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrInvalidPassword)
		repository.AssertNotCalled(t, "FindByEmail", mock.Anything, mock.Anything)
		repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("should return email already in use when user exists", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		request := dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "abc12345",
		}
		existingUser := makeAuthServiceUser(t, 1, "John Doe", "john@example.com", "abc12345")

		repository.On("FindByEmail", ctx, "john@example.com").
			Return(existingUser, nil).
			Once()

		response, err := service.Register(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrEmailAlreadyInUse)
		repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when finding user fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		expectedErr := errors.New("database unavailable")
		request := dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return((*models.User)(nil), expectedErr).
			Once()

		response, err := service.Register(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		repository.AssertExpectations(t)
	})

	t.Run("should return email already in use when create returns unique email error", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		request := dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return((*models.User)(nil), gorm.ErrRecordNotFound).
			Once()
		repository.On("Create", ctx, mock.AnythingOfType("*models.User")).
			Return(errors.New("UNIQUE constraint failed: users.email")).
			Once()

		response, err := service.Register(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrEmailAlreadyInUse)
		repository.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("should login user and return valid token", func(t *testing.T) {
		setupAuthServiceJWTEnv(t)
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		user := makeAuthServiceUser(t, 1, "John Doe", "john@example.com", "abc12345")
		request := dto_auth.LoginRequestDto{
			Email:    "  JOHN@EXAMPLE.COM  ",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return(user, nil).
			Once()

		response, err := service.Login(ctx, request)

		require.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "Bearer", response.TokenType)
		assert.Equal(t, int64(7200), response.ExpiresIn)
		assert.Equal(t, user.ID, response.User.ID)
		assert.Equal(t, user.Name, response.User.Name)
		assert.Equal(t, user.Email, response.User.Email)
		assert.Equal(t, user.Role, response.User.Role)

		claims, err := security.ValidateToken(response.Token)
		require.NoError(t, err)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
		assert.Equal(t, user.Role, claims.Role)
		repository.AssertExpectations(t)
	})

	t.Run("should return invalid credentials when user is not found", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		request := dto_auth.LoginRequestDto{
			Email:    "missing@example.com",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "missing@example.com").
			Return((*models.User)(nil), gorm.ErrRecordNotFound).
			Once()

		response, err := service.Login(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrInvalidCredentials)
		repository.AssertExpectations(t)
	})

	t.Run("should return invalid credentials when password does not match", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		user := makeAuthServiceUser(t, 1, "John Doe", "john@example.com", "abc12345")
		request := dto_auth.LoginRequestDto{
			Email:    "john@example.com",
			Password: "wrong12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return(user, nil).
			Once()

		response, err := service.Login(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrInvalidCredentials)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when finding user fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		expectedErr := errors.New("database unavailable")
		request := dto_auth.LoginRequestDto{
			Email:    "john@example.com",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return((*models.User)(nil), expectedErr).
			Once()

		response, err := service.Login(ctx, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})

	t.Run("should return error when jwt secret is empty", func(t *testing.T) {
		setupAuthServiceJWTEnv(t)
		env.JWTSecret = " "
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		user := makeAuthServiceUser(t, 1, "John Doe", "john@example.com", "abc12345")
		request := dto_auth.LoginRequestDto{
			Email:    "john@example.com",
			Password: "abc12345",
		}

		repository.On("FindByEmail", ctx, "john@example.com").
			Return(user, nil).
			Once()

		response, err := service.Login(ctx, request)

		assert.Empty(t, response)
		assert.EqualError(t, err, "JWT_SECRET não informado")
		repository.AssertExpectations(t)
	})
}

func TestMe(t *testing.T) {
	t.Run("should return authenticated user", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		user := makeAuthServiceUser(t, 1, "John Doe", "john@example.com", "abc12345")

		repository.On("FindByID", ctx, uint(1)).
			Return(user, nil).
			Once()

		response, err := service.Me(ctx, 1)

		require.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.Email, response.Email)
		assert.Equal(t, user.Role, response.Role)
		assert.Equal(t, user.CreatedAt, response.CreatedAt)
		assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
		repository.AssertExpectations(t)
	})

	t.Run("should return user not found when user does not exist", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)

		repository.On("FindByID", ctx, uint(999)).
			Return((*models.User)(nil), gorm.ErrRecordNotFound).
			Once()

		response, err := service.Me(ctx, 999)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrUserNotFound)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when finding user by id fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.UserRepositoryMock)
		service := service_auth.NewAuthService(repository)
		expectedErr := errors.New("database unavailable")

		repository.On("FindByID", ctx, uint(1)).
			Return((*models.User)(nil), expectedErr).
			Once()

		response, err := service.Me(ctx, 1)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})
}

func makeAuthServiceUser(t *testing.T, id uint, name string, email string, password string) *models.User {
	t.Helper()

	passwordHash, err := security.Hash(password)
	require.NoError(t, err)

	return &models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         shared_constants.RoleUser,
		CreatedAt:    time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC),
	}
}

func setupAuthServiceJWTEnv(t *testing.T) {
	t.Helper()

	originalJWTSecret := env.JWTSecret
	originalAppName := env.AppName
	originalExpiration := env.AuthTokenExpirationTimeInHours

	env.JWTSecret = "test-secret"
	env.AppName = "streaming_api_test"
	env.AuthTokenExpirationTimeInHours = 2

	t.Cleanup(func() {
		env.JWTSecret = originalJWTSecret
		env.AppName = originalAppName
		env.AuthTokenExpirationTimeInHours = originalExpiration
	})
}
