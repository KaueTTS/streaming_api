package api

import (
	"fmt"
	"strings"

	route_auth "github.com/KaueTTS/streaming_api/src/api/routes/auth"
	route_health "github.com/KaueTTS/streaming_api/src/api/routes/health"
	route_profile "github.com/KaueTTS/streaming_api/src/api/routes/profile"
	route_swagger "github.com/KaueTTS/streaming_api/src/api/routes/swagger"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	container "github.com/KaueTTS/streaming_api/src/container"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	redis_storage "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Init(db *gorm.DB, redisClient *redis.Client) error {
	app := fiber.New(fiber.Config{
		AppName: env.AppName,
	})

	app.Use(recover.New())
	app.Use(otelfiber.Middleware(otelfiber.WithNext(func(c *fiber.Ctx) bool {
		path := c.Path()

		return path == "/health" ||
			path == "/swagger" ||
			strings.HasPrefix(path, "/swagger/")
	})))
	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "OPTIONS,GET,PUT,DELETE,POST,PATCH",
		AllowHeaders: "Authorization,Content-Type",
	}))

	applicationContainer := container.Build(db, redisClient)

	var limiterStorage fiber.Storage
	if redisClient != nil {
		limiterStorage = redis_storage.NewFromConnection(redisClient)
	}

	injectRoutes(app, applicationContainer, limiterStorage)

	port := fmt.Sprintf(":%s", env.Port)
	if err := app.Listen(port); err != nil {
		return fmt.Errorf("falha ao iniciar o servidor: %v", err)
	}

	return nil
}

func injectRoutes(app *fiber.App, applicationContainer *container.Container, limiterStorage fiber.Storage) {
	route_health.Init(app)
	route_swagger.Init(app)

	route_auth.Init(app, applicationContainer.AuthController, limiterStorage)
	route_profile.Init(app, applicationContainer.ProfileController)
}
