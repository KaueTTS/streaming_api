package route_swagger

import (
	env "github.com/KaueTTS/streaming_api/src/configs/env"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Init(app *fiber.App) {
	if env.AppEnv == "development" {
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}
}
