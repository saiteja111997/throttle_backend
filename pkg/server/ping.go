package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

func (s *Server) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(map[string]interface{}{
		"response": "Pong",
		"status":   "Success",
	})
}

func (s *Server) GetDbData(c *fiber.Ctx) error {
	db := s.Db

	var errorData []structures.User
	db.Find(&errorData)

	return c.JSON(map[string]interface{}{
		"response": errorData,
		"status":   "Success",
	})

}

func (s *Server) UploadDbData(c *fiber.Ctx) error {
	db := s.Db

	title := c.FormValue("text")
	user_id := c.FormValue("id")

	user_id_int, err := strconv.Atoi(user_id)

	if err != nil {
		log.Fatalf("Conversion to int failed : %v", err.Error())
	}

	session_id := uuid.New().String()

	errorData := structures.Error{
		UserID: user_id_int,
		Title:  title,
		ID:     session_id,
	}

	transactionError := db.Create(&errorData).Error

	fmt.Println("Error is : ", transactionError)

	if transactionError != nil {
		return c.JSON(map[string]interface{}{
			"resposne": transactionError.Error(),
			"status":   "failed",
		})
	}

	return c.JSON(map[string]interface{}{
		"resposne": "Data inserted successfully",
		"status":   "success",
	})

}
