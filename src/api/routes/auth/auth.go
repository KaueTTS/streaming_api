package route_auth

import (
	"time"

	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// Init inicializa as rotas de autenticação.
// Configura o limiter, injeta as rotas de autenticação e as rotas privadas.
func Init(app *fiber.App, authController *v1_controller_auth.AuthController, limiterStorage fiber.Storage) {
	authLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
		Storage:    limiterStorage,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return "streaming-api:auth:" + ctx.IP()
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return responses.TooManyRequests(ctx, shared_errors_auth.TooManyAuthAttempts)
		},
	})

	authGroup := app.Group("/v1/auth", authLimiter)
	authGroup.Post("/register", authController.Register)
	authGroup.Post("/login", authController.Login)

	privateGroup := app.Group("/v1/auth", auth_middleware.AuthRequired())
	privateGroup.Get("/me", authController.Me)
}
