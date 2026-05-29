package shared_errors

import (
	"errors"

	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
)

var (
	ErrEmailAlreadyInUse  = errors.New("e-mail já está em uso")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrInvalidPassword    = errors.New("senha inválida")
	ErrUserNotFound       = errors.New("usuário não encontrado")

	ErrPasswordMustLeast8Character   = errors.New(shared_errors.PasswordMustLeast8Character)
	ErrPasswordMustMaximum72Bytes    = errors.New(shared_errors.PasswordMustMaximum72Bytes)
	ErrPasswordMustLettersAndNumbers = errors.New(shared_errors.PasswordMustLettersAndNumbers)
)
