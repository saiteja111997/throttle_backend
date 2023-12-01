package server

import "github.com/gofiber/fiber/v2"

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"response": "Pong",
		"status":   "Success",
	})
}
