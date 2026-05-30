package container

import (
	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	repository_sqlite_user "github.com/KaueTTS/streaming_api/src/repositories/sqlite/user"
	service_auth "github.com/KaueTTS/streaming_api/src/services/auth"
	"gorm.io/gorm"
)

func buildAuthController(db *gorm.DB) *v1_controller_auth.AuthController {
	userRepository := repository_sqlite_user.NewUserRepository(db)
	authService := service_auth.NewAuthService(userRepository)

	return v1_controller_auth.NewAuthController(authService)
}
