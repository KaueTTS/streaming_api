package controllers_helpers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAuthenticatedUserId(ctx *fiber.Ctx) (uint, bool) {
	userID, ok := ctx.Locals("user_id").(uint)
	return userID, ok && userID > 0
}

func ParseUintParam(ctx *fiber.Ctx, param string) (uint, string, bool) {
	value := ctx.Params(param)

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, value, false
	}

	return uint(id), value, true
}
