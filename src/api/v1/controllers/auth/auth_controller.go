package v1_controller_auth

import (
	"errors"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	validator_auth "github.com/KaueTTS/streaming_api/src/api/v1/validators"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthServiceInterface service_interface.AuthServiceInterface
}

func NewAuthController(authServiceInterface service_interface.AuthServiceInterface) *AuthController {
	return &AuthController{
		AuthServiceInterface: authServiceInterface,
	}
}

// Register godoc
// @Summary Registra um novo usuário
// @Description Cria uma nova conta de usuário no sistema
// @Tags authorization
// @Param request body dto_auth.RegisterRequestDto true "Dados para registrar o usuário"
// @Success 201 {object} dto_auth.UserResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 409 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/auth/register [post]
func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var request dto_auth.RegisterRequestDto
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

	if errDetails := validator_auth.ValidateRegisterRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_auth.InvalidRegisterData,
			errDetails,
		)
	}

	response, err := c.AuthServiceInterface.Register(ctx.Context(), request)
	if err != nil {
		if errors.Is(err, shared_errors.ErrEmailAlreadyInUse) {
			return responses.Conflict(ctx, shared_errors_auth.EmailAlreadyInUse)
		}

		if errors.Is(err, shared_errors.ErrInvalidPassword) {
			return responses.BadRequest(
				ctx,
				shared_errors_auth.InvalidPassword,
				[]dto_shared.DetailErrorDto{
					{
						Field:   shared_constants.Password,
						Value:   shared_constants.Hidden,
						Message: shared_errors_auth.InvalidPassword,
					},
				},
			)
		}

		return responses.InternalServerError(ctx, shared_errors_auth.FailedToRegisterUser)
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

// Login godoc
// @Summary Realiza login do usuário
// @Description Autentica um usuário com e-mail e senha e retorna o token de acesso
// @Tags authorization
// @Param request body dto_auth.LoginRequestDto true "Dados para autenticação"
// @Success 200 {object} dto_auth.AuthResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var request dto_auth.LoginRequestDto
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

	if errDetails := validator_auth.ValidateLoginRequest(request); len(errDetails) > 0 {
		return responses.BadRequest(
			ctx,
			shared_errors_auth.InvalidLoginData,
			errDetails,
		)
	}

	response, err := c.AuthServiceInterface.Login(ctx.Context(), request)
	if err != nil {
		if errors.Is(err, shared_errors.ErrInvalidCredentials) {
			return responses.Unauthorized(ctx, shared_errors_auth.InvalidCredentials)
		}

		return responses.InternalServerError(ctx, shared_errors_auth.FailedToLogin)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// Me godoc
// @Summary Busca o usuário autenticado
// @Description Retorna os dados do usuário autenticado com base no token enviado no header Authorization
// @Tags authorization
// @Success 200 {object} dto_auth.UserResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 404 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/auth/me [get]
// @Security BearerAuth
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return responses.Unauthorized(ctx, shared_errors_auth.InvalidToken)
	}

	response, err := c.AuthServiceInterface.Me(ctx.Context(), userID)
	if err != nil {
		if errors.Is(err, shared_errors.ErrUserNotFound) {
			return responses.NotFound(ctx, shared_errors_auth.UserNotFound)
		}

		return responses.InternalServerError(ctx, shared_errors_auth.FailedToGetAuthenticatedUser)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
