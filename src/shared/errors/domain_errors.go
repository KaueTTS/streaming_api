package shared_errors

import (
	"errors"

	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_favorite "github.com/KaueTTS/streaming_api/src/shared/errors/favorite"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
)

// Erros de domínio geral
var (
	ErrEmailAlreadyInUse     = errors.New(shared_errors_auth.EmailAlreadyInUse)
	ErrInvalidCredentials    = errors.New(shared_errors_auth.ErrInvalidCredentials)
	ErrInvalidPassword       = errors.New(shared_errors_auth.InvalidPassword)
	ErrProfileLimitReached   = errors.New(shared_errors_profile.ProfileLimitReached)
	ErrProfileNotFound       = errors.New(shared_errors_profile.ProfileNotFound)
	ErrFavoriteAlreadyExists = errors.New(shared_errors_favorite.FavoriteAlreadyExists)
	ErrUserNotFound          = errors.New(shared_errors_auth.UserNotFound)
)

// Erros de domínio relacionados a validação de senha
var (
	ErrPasswordMustLeast8Character   = errors.New(shared_errors_auth.PasswordMustLeast8Character)
	ErrPasswordMustMaximum72Bytes    = errors.New(shared_errors_auth.PasswordMustMaximum72Bytes)
	ErrPasswordMustLettersAndNumbers = errors.New(shared_errors_auth.PasswordMustLettersAndNumbers)
)
