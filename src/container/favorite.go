package container

import (
	v1_controller_favorite "github.com/KaueTTS/streaming_api/src/api/v1/controllers/favorite"
	repository_sqlite_favorite "github.com/KaueTTS/streaming_api/src/repositories/sqlite/favorite"
	service_favorite "github.com/KaueTTS/streaming_api/src/services/favorite"
	"gorm.io/gorm"
)

func buildFavoriteController(db *gorm.DB) *v1_controller_favorite.FavoriteController {
	favoriteRepository := repository_sqlite_favorite.NewFavoriteRepository(db)
	favoriteService := service_favorite.NewFavoriteService(favoriteRepository)

	return v1_controller_favorite.NewFavoriteController(favoriteService)
}
