package route_content

import (
	v1_controller_content "github.com/KaueTTS/streaming_api/src/api/v1/controllers/content"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, contentController *v1_controller_content.ContentController) {
	privateGroup := app.Group("/v1", auth_middleware.AuthRequired())
	privateGroup.Get("/contents", contentController.ListContents)
	privateGroup.Get("/contents/search", contentController.SearchContents)
}
