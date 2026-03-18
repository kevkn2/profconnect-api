package handler

import "github.com/gofiber/fiber/v3"

func (h *Handler) HelloWorld(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello World!",
	})
}
