package container

import (
	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	AuthController    *v1_controller_auth.AuthController
	ProfileController *v1_controller_profile.ProfileController
	RedisClient       *redis.Client
}

func Build(db *gorm.DB, redisClient *redis.Client) *Container {
	return &Container{
		AuthController:    buildAuthController(db),
		ProfileController: buildProfileController(db),
		RedisClient:       redisClient,
	}
}
