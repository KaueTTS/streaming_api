package v1_controller_profile

import (
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	service_interface "github.com/KaueTTS/streaming_api/src/services/interfaces"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
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

// Profiles godoc
// @Summary Teste
// @Description Teste
// @Tags profiles
// @Success 200 {object} []dto_profile.ProfileResponseDto
// @Failure 400 {object} dto_shared.ErrorDto
// @Failure 401 {object} dto_shared.ErrorDto
// @Failure 500 {object} dto_shared.ErrorDto
// @Router /v1/profiles [get]
// @Security BearerAuth
func (c *ProfileController) Profiles(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(uint)
	if !ok || userID == 0 {
		return responses.Unauthorized(ctx, shared_errors.InvalidToken)
	}

	role, _ := ctx.Locals("user_role").(string)

	profiles, err := c.ProfileServiceInterface.GetProfiles(ctx.Context(), userID, role)
	if err != nil {
		return responses.InternalServerError(ctx, shared_errors.FailedToListProfiles)
	}

	return ctx.Status(fiber.StatusOK).JSON(profiles)
}
