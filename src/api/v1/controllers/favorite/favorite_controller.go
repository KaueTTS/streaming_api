package v1_controller_favorite

import (
	"errors"
	"fmt"

	controllers_helpers "github.com/KaueTTS/streaming_api/src/api/v1/controllers"
	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_favorite "github.com/KaueTTS/streaming_api/src/shared/errors/favorite"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
	"github.com/gofiber/fiber/v2"
)

type FavoriteController struct {
	favoriteService service_interface.FavoriteServiceInterface
}

func NewFavoriteController(favoriteService service_interface.FavoriteServiceInterface) *FavoriteController {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// ListFavorites godoc
// @Summary Lista os favoritos de um perfil
// @Description Retorna uma lista paginada de filmes e séries favoritos de um perfil específico.
// @Tags favorites
// @Param profile_id query int true "ID do perfil"
// @Param page query int false "Número da página" default(1)
// @Param per_page query int false "Número de itens por página" default(10)
// @Success 200 {object} dto_favorite.FavoriteResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 404 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/favorites [get]
// @Security BearerAuth
func (c *FavoriteController) ListFavorites(ctx *fiber.Ctx) error {
	userID, ok := controllers_helpers.GetAuthenticatedUserID(ctx)
	if !ok {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	request, details, ok := controllers_helpers.ParseQuery[dto_favorite.FavoriteListRequestDto](ctx)
	if !ok {
		return responses.BadRequest(ctx, shared_errors.InvalidQueryParameters, details)
	}

	if request.ProfileID == 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_favorite.InvalidProfileID,
			[]dto_shared.DetailErrorDto{
				{
					Field:   shared_constants.ProfileID,
					Value:   ctx.Query(shared_constants.ProfileID),
					Message: shared_errors_favorite.InvalidProfileID,
				},
			},
		)
	}

	page, perPage := shared_normalizers.NormalizePagination(dto_shared.PaginationDto{
		Page:    request.Page,
		PerPage: request.PerPage,
	})

	response, err := c.favoriteService.ListFavorites(ctx.UserContext(), userID, request.ProfileID, page, perPage)
	if err != nil {
		if errors.Is(err, shared_errors.ErrProfileNotFound) {
			return responses.NotFound(ctx, shared_errors_profile.ProfileNotFound)
		}

		return responses.InternalServerError(ctx, shared_errors_favorite.FailedToListFavorites)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// AddFavorite godoc
// @Summary
// @Description
// @Tags favorites
// @Success 200
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/favorites [post]
// @Security BearerAuth
func (c *FavoriteController) AddFavorite(ctx *fiber.Ctx) error {
	return fmt.Errorf("não implementado")
}

// RemoveFavorite godoc
// @Summary
// @Description
// @Tags favorites
// @Success 200
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/favorites/{id} [delete]
// @Security BearerAuth
func (c *FavoriteController) RemoveFavorite(ctx *fiber.Ctx) error {
	return fmt.Errorf("não implementado")
}
