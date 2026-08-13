package api

import (
	"fmt"
	"strings"

	route_auth "github.com/KaueTTS/streaming_api/src/api/routes/auth"
	route_content "github.com/KaueTTS/streaming_api/src/api/routes/content"
	route_favorite "github.com/KaueTTS/streaming_api/src/api/routes/favorite"
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

// Init inicia a aplicação e configura as rotas.
// Configura o recover, otel, helmet, cors, injeta as rotas e o container.
// Usa o redisClient para criar o limiterStorage se o redisClient for diferente de nil.
// Retorna uma instância da aplicação.
func Init(db *gorm.DB, redisClient *redis.Client) *fiber.App {
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

	appContainer := container.Build(db, redisClient)

	var limiterStorage fiber.Storage
	if redisClient != nil {
		limiterStorage = redis_storage.NewFromConnection(redisClient)
	}

	injectRoutes(app, appContainer, limiterStorage)

	return app
}

// Listen inicia o servidor na porta definida em env.Port
func Listen(app *fiber.App) error {
	return app.Listen(fmt.Sprintf(":%s", env.Port))
}

// injectRoutes injeta as rotas no aplicativo.
func injectRoutes(app *fiber.App, container *container.Container, limiterStorage fiber.Storage) {
	route_health.Init(app)
	route_swagger.Init(app)

	route_auth.Init(app, container.AuthController, limiterStorage)
	route_content.Init(app, container.ContentController)
	route_profile.Init(app, container.ProfileController)
	route_favorite.Init(app, container.FavoriteController)
}
