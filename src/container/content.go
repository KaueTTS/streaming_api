package container

import (
	v1_controller_content "github.com/KaueTTS/streaming_api/src/api/v1/controllers/content"
	repository_http_tmdb "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb"
	service_content "github.com/KaueTTS/streaming_api/src/services/content"
)

func buildContentController() *v1_controller_content.ContentController {
	tmdbRepository := repository_http_tmdb.NewTMDBRepository()
	contentService := service_content.NewContentService(tmdbRepository)

	return v1_controller_content.NewContentController(contentService)
}
