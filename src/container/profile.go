package container

import (
	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	repository_sqlite_profile "github.com/KaueTTS/streaming_api/src/repositories/sqlite/profile"
	service_profile "github.com/KaueTTS/streaming_api/src/services/profile"
	"gorm.io/gorm"
)

// buildProfileController constrói o ProfileController com as dependências necessárias
func buildProfileController(db *gorm.DB) *v1_controller_profile.ProfileController {
	profileRepository := repository_sqlite_profile.NewProfileRepository(db)
	profileService := service_profile.NewProfileService(profileRepository)

	return v1_controller_profile.NewProfileController(profileService)
}
