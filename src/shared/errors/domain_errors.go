package shared_errors

import (
	"errors"

	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
)

var (
	ErrEmailAlreadyInUse   = errors.New(shared_errors_auth.EmailAlreadyInUse)
	ErrInvalidCredentials  = errors.New("credenciais inválidas")
	ErrInvalidPassword     = errors.New(shared_errors_auth.InvalidPassword)
	ErrProfileLimitReached = errors.New("limite de perfis por usuário atingido")
	ErrUserNotFound        = errors.New(shared_errors_auth.UserNotFound)

	ErrPasswordMustLeast8Character   = errors.New(shared_errors_auth.PasswordMustLeast8Character)
	ErrPasswordMustMaximum72Bytes    = errors.New(shared_errors_auth.PasswordMustMaximum72Bytes)
	ErrPasswordMustLettersAndNumbers = errors.New(shared_errors_auth.PasswordMustLettersAndNumbers)
)
