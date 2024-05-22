package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Register(c *fiber.Ctx) error {

	fullName := c.FormValue("full_name")
	password := c.FormValue("password")
	email := c.FormValue("email")

	fmt.Println("Printing input values : ", fullName, password, email)

	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := structures.Users{
		FullName: fullName,
		Email:    email,
		Password: hashedPassword,
	}

	if err := s.Db.Create(&user).Error; err != nil {
		log.Fatal("Error is : ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})

}

func (s *Server) Login(c *fiber.Ctx) error {

	password := c.FormValue("password")
	email := c.FormValue("email")

	var user structures.Users

	result := s.Db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Login successful"})

}
