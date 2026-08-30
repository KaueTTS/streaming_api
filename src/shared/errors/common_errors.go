package shared_errors

// Erros comuns para campos de requisição e consulta
const (
	InvalidRequestBody     = "corpo da requisição inválido"
	InvalidQueryParameters = "parâmetros de consulta inválidos"
)

// Erros comuns para campo name
const (
	NameRequired                = "nome é obrigatório"
	NameMustLeast2Character     = "nome deve ter pelo menos 2 caracteres"
	NameMustMaximum120Character = "nome deve ter no máximo 120 caracteres"
)

// Erros comuns para campo user id
const (
	InvalidUserID = "id do usuário inválido"
)

// Erros comuns para campo access
const (
	AccessAdminOnly = "acesso permitido somente para admin"
)

// Erros comuns para campos de paginação
const (
	PageMustBePositive    = "page deve ser um número inteiro positivo"
	PerPageMustBePositive = "per_page deve ser um número inteiro positivo"
)
