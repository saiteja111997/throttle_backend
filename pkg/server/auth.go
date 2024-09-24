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

	userName := c.FormValue("user_name")
	password := c.FormValue("password")

	fmt.Println("Printing input values: ", userName, password)

	// Check if the username already exists
	var existingUser structures.Users
	if err := s.Db.Where("username = ?", userName).First(&existingUser).Error; err == nil {
		// Username exists, return a conflict response
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	// Hash the password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create new user record
	user := structures.Users{
		Username: userName,
		Password: hashedPassword,
	}

	if err := s.Db.Create(&user).Error; err != nil {
		log.Fatal("Error is: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	fmt.Println("Printing user id :", user.ID)

	// Return the newly created user's ID
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"userId":  user.ID, // Assuming user.ID is the auto-generated ID field
	})
}

func (s *Server) Login(c *fiber.Ctx) error {

	userName := c.FormValue("user_name")
	password := c.FormValue("password")

	var user structures.Users

	result := s.Db.Where("username = ?", userName).First(&user)
	if result.Error != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication failed"})
	}

	// Update the signed in field in the database to 1 for the logged in user
	// err := s.Db.Raw("UPDATE users SET signed_in =? WHERE id =?", "1", user.ID).Error
	// if err != nil {
	// 	c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to update user login status"})
	// }

	// fmt.Println("Printing user id : ", user.ID)
	// fmt.Println("Printing user signed in status : ", user.SignedIn)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"userId":  user.ID,
	})

}

// func (s *Server) Logout(c *fiber.Ctx) error {
// 	user_id := c.FormValue("user_id")

// 	// Update the signed in field in the database to 0 for the logged in user
// 	err := s.Db.Raw("UPDATE users SET signed_in = ? WHERE id = ?", "0", user_id).Error
// 	if err != nil {
// 		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to logout"})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Action successful"})
// }

// func (s *Server) GetAuthStatus(c *fiber.Ctx) error {
// 	user_id := c.FormValue("user_id")

// 	// Get signed in status from db
// 	var user structures.Users

// 	result := s.Db.Where("id =?", user_id).First(&user)
// 	if result.Error != nil {
// 		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"signed_in": user.SignedIn,
// 	})
// }
