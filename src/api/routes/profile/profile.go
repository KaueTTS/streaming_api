package route_profile

import (
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	repository_sqlite_profile "github.com/KaueTTS/streaming_api/src/repositories/sqlite/profile"
	service_profile "github.com/KaueTTS/streaming_api/src/services/profile"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB) {
	profileRepository := repository_sqlite_profile.NewProfileRepository(db)
	profileService := service_profile.NewProfileService(profileRepository)
	profileController := v1_controller_profile.NewProfileController(profileService)

	privateGroup := app.Group("/v1", auth_middleware.AuthRequired())
	privateGroup.Get("/profiles", profileController.ListProfiles)
	privateGroup.Post("/profiles", profileController.CreateProfile)
}
