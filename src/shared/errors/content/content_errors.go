package shared_errors_content

const (
	InvalidContentQueryParameters = "parâmetros de consulta de conteúdo inválidos"
	InvalidContentType            = "type deve ser movie ou tv"
	InvalidContentSortBy          = "sort_by inválido para o type informado"
	InvalidContentGenres          = "with_genres deve conter apenas ids numéricos separados por vírgula ou pipe"
	InvalidContentLanguage        = "language deve estar no formato ISO, exemplo: pt-BR"
	InvalidContentYear            = "year deve ser maior que zero quando informado"
)

const (
	ContentTypeRequired = "type é obrigatório"
	SearchQueryRequired = "query é obrigatório"
)

const (
	FailedToListContents   = "falha ao listar conteúdos"
	FailedToSearchContents = "falha ao buscar conteúdos"
)
