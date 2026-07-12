package container

import (
	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	v1_controller_content "github.com/KaueTTS/streaming_api/src/api/v1/controllers/content"
	v1_controller_favorite "github.com/KaueTTS/streaming_api/src/api/v1/controllers/favorite"
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	AuthController     *v1_controller_auth.AuthController
	ContentController  *v1_controller_content.ContentController
	ProfileController  *v1_controller_profile.ProfileController
	FavoriteController *v1_controller_favorite.FavoriteController
	RedisClient        *redis.Client
}

func Build(db *gorm.DB, redisClient *redis.Client) *Container {
	return &Container{
		AuthController:     buildAuthController(db),
		ContentController:  buildContentController(),
		FavoriteController: buildFavoriteController(db),
		ProfileController:  buildProfileController(db),
		RedisClient:        redisClient,
	}
}
