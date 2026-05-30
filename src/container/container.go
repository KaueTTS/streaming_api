package container

import (
	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	"gorm.io/gorm"
)

type Container struct {
	AuthController    *v1_controller_auth.AuthController
	ProfileController *v1_controller_profile.ProfileController
}

func Build(db *gorm.DB) *Container {
	return &Container{
		AuthController:    buildAuthController(db),
		ProfileController: buildProfileController(db),
	}
}
