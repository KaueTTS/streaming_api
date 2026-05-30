package v1_controller_profile

import (
	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	validator_profile "github.com/KaueTTS/streaming_api/src/api/v1/validators"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	ProfileServiceInterface service_interface.ProfileServiceInterface
}

func NewProfileController(profileServiceInterface service_interface.ProfileServiceInterface) *ProfileController {
	return &ProfileController{
		ProfileServiceInterface: profileServiceInterface,
	}
}

// ListProfiles godoc
// @Summary Listar os perfis do usuário logado
// @Description Retorna uma lista paginada dos perfis associados ao usuário autenticado. O usuário pode ter múltiplos perfis, e esta rota permite recuperar todos eles de forma organizada e eficiente.
// @Param page query int false "Número da página" default(1)
// @Param per_page query int false "Número de itens por página" default(10)
// @Tags profiles
// @Success 200 {object} dto_profile.ProfileResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/profiles [get]
// @Security BearerAuth
func (c *ProfileController) ListProfiles(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	var pagination dto_shared.PaginationDto
	if err := ctx.QueryParser(&pagination); err != nil {
		return responses.BadRequest(
			ctx,
			shared_errors.InvalidQueryParameters,
			[]dto_shared.DetailErrorDto{
				{
					Field:   shared_constants.Page,
					Value:   ctx.Query(shared_constants.Page),
					Message: shared_errors.PageMustBePositive,
				},
				{
					Field:   shared_constants.PerPage,
					Value:   ctx.Query(shared_constants.PerPage),
					Message: shared_errors.PerPageMustBePositive,
				},
			},
		)
	}

	page, perPage := validator_profile.ValidatePagination(pagination)

	response, err := c.ProfileServiceInterface.ListProfiles(ctx.UserContext(), userID, page, perPage)
	if err != nil {
		return responses.InternalServerError(ctx, shared_errors_profile.FailedToListProfiles)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// CreateProfile godoc
// @Summary Criar um novo perfil para o usuário logado
// @Description Cria um novo perfil associado ao usuário autenticado. O usuário pode ter múltiplos perfis, e esta rota permite criar um novo perfil com as informações fornecidas no corpo da requisição.
// @Param request body dto_profile.CreateProfileRequestDto true "Dados para criar um perfil"
// @Tags profiles
// @Success 201 {object} dto_profile.ProfileDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/profiles [post]
// @Security BearerAuth
func (c *ProfileController) CreateProfile(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	var request dto_profile.CreateProfileRequestDto
	if err := ctx.BodyParser(&request); err != nil {
		return responses.BadRequest(
			ctx,
			shared_errors.InvalidRequestBody,
			[]dto_shared.DetailErrorDto{
				{
					Field:   "",
					Value:   "",
					Message: err.Error(),
				},
			},
		)
	}

	if errDetails := validator_profile.ValidateCreateProfileRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_profile.InvalidCreateProfileData,
			errDetails,
		)
	}

	response, err := c.ProfileServiceInterface.CreateProfile(ctx.UserContext(), userID, request)
	if err != nil {
		return responses.InternalServerError(ctx, shared_errors_profile.FailedToCreateProfile)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (c *ProfileController) UpdateProfile(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"message": "not implemented"})
}

func (c *ProfileController) DeleteProfile(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"message": "not implemented"})
}
