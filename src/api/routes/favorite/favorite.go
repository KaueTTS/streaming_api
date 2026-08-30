package route_favorite

import (
	v1_controller_favorite "github.com/KaueTTS/streaming_api/src/api/v1/controllers/favorite"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, favoriteController *v1_controller_favorite.FavoriteController) {
	privateGroup := app.Group("/v1", auth_middleware.AuthRequired())
	privateGroup.Get("/favorites", favoriteController.ListFavorites)
	privateGroup.Post("/favorites", favoriteController.AddFavorite)
	privateGroup.Delete("/favorites", favoriteController.RemoveFavorite)
}
