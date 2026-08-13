package route_health

import "github.com/gofiber/fiber/v2"

type Health struct {
	Status string `json:"status"`
}

// healthRoute 	godoc
// @Summary 	Verifica se a API está online
// @Description Retorna o status de funcionamento da API
// @Tags 		health
// @Success 	200 {object} Health
// @Router 		/health [get]
func healthRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(Health{
		Status: "ok",
	})
}

func Init(app *fiber.App) {
	app.Get("/health", healthRoute)
}
