package container

import (
	v1_controller_content "github.com/KaueTTS/streaming_api/src/api/v1/controllers/content"
	repository_http_tmdb "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb"
	repository_sqlite_profile "github.com/KaueTTS/streaming_api/src/repositories/sqlite/profile"
	service_content "github.com/KaueTTS/streaming_api/src/services/content"
	"gorm.io/gorm"
)

// buildContentController constrói o ContentController com as dependências necessárias
func buildContentController(db *gorm.DB) *v1_controller_content.ContentController {
	tmdbRepository := repository_http_tmdb.NewTMDBRepository()
	profileRepository := repository_sqlite_profile.NewProfileRepository(db)

	contentService := service_content.NewContentService(tmdbRepository, profileRepository)

	return v1_controller_content.NewContentController(contentService)
}
