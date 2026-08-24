package shared_errors_favorite

// Erros comuns de operações de favorite
const (
	FailedToListFavorites  = "erro ao listar favoritos"
	FailedToCreateFavorite = "erro ao criar favorito"
	InvalidProfileID       = "id do perfil inválido"
	FavoriteAlreadyExists  = "Conteúdo já adicionado aos favoritos"
)

// Erros de validação de favorite
const (
	InvalidCreateFavoriteData = "dados inválidos para criar favorito"
	InvalidContentExternalID  = "id externo do conteúdo inválido"
	InvalidContentType        = "tipo de conteúdo precisa ser movie ou tv"
)
