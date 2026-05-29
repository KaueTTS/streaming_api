package route_auth

import (
	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	repository_sqlite_user "github.com/KaueTTS/streaming_api/src/repositories/sqlite/user"
	service_auth "github.com/KaueTTS/streaming_api/src/services/auth"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gorm.io/gorm"
)

func Init(app *fiber.App, db *gorm.DB) {
	userRepository := repository_sqlite_user.NewUserRepository(db)
	authService := service_auth.NewAuthService(userRepository)
	authController := v1_controller_auth.NewAuthController(authService)

	authLimiter := limiter.New(limiter.Config{
		Max:        shared_constants.AuthRateLimitMax,
		Expiration: shared_constants.AuthRateLimitExpiration,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
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
