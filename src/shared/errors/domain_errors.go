package shared_errors

import (
	"errors"
)

var (
	ErrEmailAlreadyInUse  = errors.New("e-mail já está em uso")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrInvalidPassword    = errors.New("senha inválida")
	ErrUserNotFound       = errors.New("usuário não encontrado")

	ErrPasswordMustLeast8Character   = errors.New(PasswordMustLeast8Character)
	ErrPasswordMustMaximum72Bytes    = errors.New(PasswordMustMaximum72Bytes)
	ErrPasswordMustLettersAndNumbers = errors.New(PasswordMustLettersAndNumbers)
)
