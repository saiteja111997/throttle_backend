package server

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) HealthCheck(c *fiber.Ctx) error {

	for i := 0; i <= 10; i++ {
		pingErr := s.Db.DB().Ping()
		if pingErr != nil {
			c.JSON(map[string]interface{}{
				"response": "DB connection not alive",
				"status":   "failed",
			})
		} else {
			fmt.Println("DB connection is alive!!")
			break
		}

		time.Sleep(2 * time.Second)
	}

	return c.JSON(map[string]interface{}{
		"response": "Pong",
		"status":   "Success",
	})
}
