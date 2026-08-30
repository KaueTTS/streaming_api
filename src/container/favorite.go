package container

import (
	v1_controller_favorite "github.com/KaueTTS/streaming_api/src/api/v1/controllers/favorite"
	repository_http_tmdb "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb"
	repository_sqlite_favorite "github.com/KaueTTS/streaming_api/src/repositories/sqlite/favorite"
	service_favorite "github.com/KaueTTS/streaming_api/src/services/favorite"
	"gorm.io/gorm"
)

// buildFavoriteController constrói o FavoriteController com as dependências necessárias
func buildFavoriteController(db *gorm.DB) *v1_controller_favorite.FavoriteController {
	favoriteRepository := repository_sqlite_favorite.NewFavoriteRepository(db)
	tmdbRepository := repository_http_tmdb.NewTMDBRepository()

	favoriteService := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)

	return v1_controller_favorite.NewFavoriteController(favoriteService)
}
