package shared_errors

const (
	InvalidRequestBody     = "corpo da requisição inválido"
	InvalidQueryParameters = "parâmetros de consulta inválidos"
)

const (
	NameRequired                = "nome é obrigatório"
	NameMustLeast2Character     = "nome deve ter pelo menos 2 caracteres"
	NameMustMaximum120Character = "nome deve ter no máximo 120 caracteres"
)

const (
	InvalidUserID = "id do usuário inválido"
)

const (
	AccessAdminOnly = "acesso permitido somente para admin"
)

const (
	PageMustBePositive    = "page deve ser um número inteiro positivo"
	PerPageMustBePositive = "per_page deve ser um número inteiro positivo"
)
