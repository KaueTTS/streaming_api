package shared_errors

const (
	InvalidRequestBody     = "corpo da requisição inválido"
	InvalidQueryParameters = "parâmetros de consulta inválidos"

	InvalidUserID = "id do usuário inválido"

	AccessAdminOnly = "acesso permitido somente para admin"

	PageMustBePositive    = "page deve ser um número inteiro positivo"
	PerPageMustBePositive = "per_page deve ser um número inteiro positivo"
)
