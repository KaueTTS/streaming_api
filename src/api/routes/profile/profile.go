package route_profile

import (
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App, profileController *v1_controller_profile.ProfileController) {
	privateGroup := app.Group("/v1", auth_middleware.AuthRequired())
	privateGroup.Get("/profiles", profileController.ListProfiles)
	privateGroup.Post("/profiles", profileController.CreateProfile)
}
