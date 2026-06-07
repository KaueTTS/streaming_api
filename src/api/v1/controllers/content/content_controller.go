package v1_controller_content

import (
	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	validator_content "github.com/KaueTTS/streaming_api/src/api/v1/validators"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_content "github.com/KaueTTS/streaming_api/src/shared/errors/content"
	"github.com/gofiber/fiber/v2"
)

type ContentController struct {
	ContentServiceInterface service_interface.ContentServiceInterface
}

func NewContentController(contentServiceInterface service_interface.ContentServiceInterface) *ContentController {
	return &ContentController{
		ContentServiceInterface: contentServiceInterface,
	}
}

// ListContents godoc
// @Summary Listar conteúdos
// @Description Lista filmes ou séries usando a integração com a TMDB API
// @Tags contents
// @Param type query string true "Tipo do conteúdo" Enums(movie,tv)
// @Param page query int false "Número da página" default(1)
// @Param sort_by query string false "Ordenação usada pela TMDB. Exemplo: popularity.desc"
// @Param with_genres query string false "IDs de gêneros separados por vírgula"
// @Param language query string false "Idioma da resposta. Exemplo: pt-BR"
// @Param year query int false "Ano de lançamento"
// @Success 200 {object} dto_content.ContentResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/contents [get]
// @Security BearerAuth
func (c *ContentController) ListContents(ctx *fiber.Ctx) error {
	var request dto_content.ContentListRequestDto
	if err := ctx.QueryParser(&request); err != nil {
		return responses.BadRequest(
			ctx,
			shared_errors.InvalidQueryParameters,
			[]dto_shared.DetailErrorDto{
				{
					Field:   shared_constants.Page,
					Value:   ctx.Query(shared_constants.Page),
					Message: shared_errors.PageMustBePositive,
				},
			},
		)
	}

	if errDetails := validator_content.ValidateContentListRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_content.InvalidContentQueryParameters,
			errDetails,
		)
	}

	response, err := c.ContentServiceInterface.ListContents(ctx.UserContext(), request)
	if err != nil {
		return responses.InternalServerError(ctx, shared_errors_content.FailedToListContents)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// SearchContents godoc
// @Summary Buscar conteúdos
// @Description Busca filmes ou séries usando a integração com a TMDB API
// @Tags contents
// @Param type query string true "Tipo do conteúdo" Enums(movie,tv)
// @Param query query string true "Termo pesquisado"
// @Param page query int false "Número da página" default(1)
// @Param language query string false "Idioma da resposta. Exemplo: pt-BR"
// @Success 200 {object} dto_content.ContentResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/contents/search [get]
// @Security BearerAuth
func (c *ContentController) SearchContents(ctx *fiber.Ctx) error {
	var request dto_content.ContentSearchRequestDto
	if err := ctx.QueryParser(&request); err != nil {
		return responses.BadRequest(
			ctx,
			shared_errors.InvalidQueryParameters,
			[]dto_shared.DetailErrorDto{
				{
					Field:   shared_constants.Page,
					Value:   ctx.Query(shared_constants.Page),
					Message: shared_errors.PageMustBePositive,
				},
			},
		)
	}

	if errDetails := validator_content.ValidateContentSearchRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_content.InvalidContentQueryParameters,
			errDetails,
		)
	}

	response, err := c.ContentServiceInterface.SearchContents(ctx.UserContext(), request)
	if err != nil {
		return responses.InternalServerError(ctx, shared_errors_content.FailedToSearchContents)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
