package handler

import "github.com/gofiber/fiber/v3"

// HelloWorld godoc
// @Summary Hello World
// @Description Returns a simple hello world message
// @Tags Health
// @Produce json
// @Success 200 Hello World!
// @Router /hello [get]
func (h *Handler) HelloWorld(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello World!",
	})
}
