package v1_controller_content

import (
	"errors"

	controllers_helpers "github.com/KaueTTS/streaming_api/src/api/v1/controllers"
	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	validator_content "github.com/KaueTTS/streaming_api/src/api/v1/validators"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_content "github.com/KaueTTS/streaming_api/src/shared/errors/content"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	"github.com/gofiber/fiber/v2"
)

type ContentController struct {
	contentService service_interface.ContentServiceInterface
}

func NewContentController(contentService service_interface.ContentServiceInterface) *ContentController {
	return &ContentController{
		contentService: contentService,
	}
}

// ListContents godoc
// @Summary 	Listar conteúdos
// @Description Lista filmes ou séries usando a integração com a TMDB API
// @Tags 		contents
// @Param 		type query string true "Tipo do conteúdo" Enums(movie,tv)
// @Param 		profile_id query uint true "ID do perfil"
// @Param 		page query int false "Número da página" default(1)
// @Param 		sort_by query string false "Ordenação usada pela TMDB. Exemplo: popularity.desc"
// @Param 		with_genres query string false "IDs de gêneros separados por vírgula"
// @Param 		language query string false "Idioma da resposta. Exemplo: pt-BR"
// @Param 		year query int false "Ano de lançamento"
// @Success 	200 {object} dto_content.ContentResponseDto
// @Failure 	400 {object} dto_shared.ErrorDto
// @Failure 	401 {object} dto_shared.ErrorDto
// @Failure 	404 {object} dto_shared.ErrorDto
// @Failure 	502 {object} dto_shared.ErrorDto
// @Router 		/v1/contents [get]
// @Security 	BearerAuth
func (c *ContentController) ListContents(ctx *fiber.Ctx) error {
	userID, ok := controllers_helpers.GetAuthenticatedUserID(ctx)
	if !ok {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	request, details, ok := controllers_helpers.ParseQuery[dto_content.ContentListRequestDto](ctx)
	if !ok {
		return responses.BadRequest(ctx, shared_errors.InvalidQueryParameters, details)
	}

	if errDetails := validator_content.ValidateContentListRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_content.InvalidContentQueryParameters,
			errDetails,
		)
	}

	response, err := c.contentService.ListContents(ctx.UserContext(), userID, request)
	if err != nil {
		if errors.Is(err, shared_errors.ErrProfileNotFound) {
			return responses.NotFound(ctx, shared_errors_profile.ProfileNotFound)
		}

		return responses.BadGateway(ctx, shared_errors_content.FailedToListContents)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// SearchContents 	godoc
// @Summary 		Buscar conteúdos
// @Description 	Busca filmes ou séries usando a integração com a TMDB API
// @Tags 			contents
// @Param 			type query string true "Tipo do conteúdo" Enums(movie,tv)
// @Param 			profile_id query uint true "ID do perfil"
// @Param 			query query string true "Termo pesquisado"
// @Param 			page query int false "Número da página" default(1)
// @Param 			language query string false "Idioma da resposta. Exemplo: pt-BR"
// @Success 		200 {object} dto_content.ContentResponseDto
// @Failure 		400 {object} dto_shared.ErrorDto
// @Failure 		401 {object} dto_shared.ErrorDto
// @Failure 		404 {object} dto_shared.ErrorDto
// @Failure 		502 {object} dto_shared.ErrorDto
// @Router 			/v1/contents/search [get]
// @Security 		BearerAuth
func (c *ContentController) SearchContents(ctx *fiber.Ctx) error {
	userID, ok := controllers_helpers.GetAuthenticatedUserID(ctx)
	if !ok {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	request, details, ok := controllers_helpers.ParseQuery[dto_content.ContentSearchRequestDto](ctx)
	if !ok {
		return responses.BadRequest(ctx, shared_errors.InvalidQueryParameters, details)
	}

	if errDetails := validator_content.ValidateContentSearchRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_content.InvalidContentQueryParameters,
			errDetails,
		)
	}

	response, err := c.contentService.SearchContents(ctx.UserContext(), userID, request)
	if err != nil {
		if errors.Is(err, shared_errors.ErrProfileNotFound) {
			return responses.NotFound(ctx, shared_errors_profile.ProfileNotFound)
		}

		return responses.BadGateway(ctx, shared_errors_content.FailedToSearchContents)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
