package shared_errors_auth

const (
	InvalidRegisterData = "dados de cadastro inválidos"
	InvalidLoginData    = "dados de login inválidos"

	FailedToRegisterUser = "erro ao cadastrar usuário"
	FailedToLogin        = "erro ao realizar login"

	TokenMissingOrInvalid = "token não informado ou inválido"
	TokenInvalidOrExpired = "token inválido ou expirado"
	InvalidToken          = "token inválido"

	UserNotFound                 = "usuário não encontrado"
	FailedToGetAuthenticatedUser = "erro ao buscar usuário autenticado"

	TooManyAuthAttempts = "muitas tentativas, tente novamente em alguns instantes"
)

const (
	NameRequired                = "nome é obrigatório"
	NameMustLeast2Character     = "nome deve ter pelo menos 2 caracteres"
	NameMustMaximum120Character = "nome deve ter no máximo 120 caracteres"
)

const (
	PasswordRequired              = "senha é obrigatória"
	PasswordMustLeast8Character   = "senha deve ter pelo menos 8 caracteres"
	PasswordMustMaximum72Bytes    = "senha deve ter no máximo 72 bytes"
	PasswordMustLettersAndNumbers = "senha deve conter letras e números"
	InvalidPassword               = "senha inválida"
)

const (
	EmailRequired      = "e-mail é obrigatório"
	EmailInvalid       = "e-mail inválido"
	EmailAlreadyInUse  = "e-mail já está em uso"
	InvalidCredentials = "e-mail ou senha inválidos"
)
