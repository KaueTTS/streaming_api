package service_auth

import (
	"context"
	"errors"
	"strings"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	models "github.com/KaueTTS/streaming_api/src/models"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	security "github.com/KaueTTS/streaming_api/src/security"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepositoryInterface repository_interface.AuthRepositoryInterface
}

func NewAuthService(authRepositoryInterface repository_interface.AuthRepositoryInterface) *AuthService {
	return &AuthService{
		AuthRepositoryInterface: authRepositoryInterface,
	}
}

func (s *AuthService) Register(ctx context.Context, request dto_auth.RegisterRequestDto) (dto_auth.UserResponseDto, error) {
	name := strings.TrimSpace(request.Name)
	email := normalizeEmail(request.Email)

	if err := security.Validate(request.Password); err != nil {
		return dto_auth.UserResponseDto{}, shared_errors.ErrInvalidPassword
	}

	existingUser, err := s.AuthRepositoryInterface.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_auth.UserResponseDto{}, err
	}

	if existingUser != nil {
		return dto_auth.UserResponseDto{}, shared_errors.ErrEmailAlreadyInUse
	}

	passwordHash, err := security.Hash(request.Password)
	if err != nil {
		return dto_auth.UserResponseDto{}, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	if err := s.AuthRepositoryInterface.Create(ctx, user); err != nil {
		if isUniqueConstraintError(err) {
			return dto_auth.UserResponseDto{}, shared_errors.ErrEmailAlreadyInUse
		}

		return dto_auth.UserResponseDto{}, err
	}

	return toUserResponseDto(*user), nil
}

func (s *AuthService) Login(ctx context.Context, request dto_auth.LoginRequestDto) (dto_auth.AuthResponseDto, error) {
	email := normalizeEmail(request.Email)

	user, err := s.AuthRepositoryInterface.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_auth.AuthResponseDto{}, shared_errors.ErrInvalidCredentials
		}

		return dto_auth.AuthResponseDto{}, err
	}

	if !security.Compare(user.PasswordHash, request.Password) {
		return dto_auth.AuthResponseDto{}, shared_errors.ErrInvalidCredentials
	}

	token, err := security.GenerateToken(user.ID, user.Email)
	if err != nil {
		return dto_auth.AuthResponseDto{}, err
	}

	return dto_auth.AuthResponseDto{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(security.GetExpirationDuration().Seconds()),
		User:      toUserResponseDto(*user),
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (dto_auth.UserResponseDto, error) {
	user, err := s.AuthRepositoryInterface.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_auth.UserResponseDto{}, shared_errors.ErrUserNotFound
		}

		return dto_auth.UserResponseDto{}, err
	}

	return toUserResponseDto(*user), nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func toUserResponseDto(user models.User) dto_auth.UserResponseDto {
	return dto_auth.UserResponseDto{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func isUniqueConstraintError(err error) bool {
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "unique") && strings.Contains(message, "email")
}
