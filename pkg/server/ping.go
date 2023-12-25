package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"response": "Pong",
		"status":   "Success",
	})
}
