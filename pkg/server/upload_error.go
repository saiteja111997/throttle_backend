package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
)

const (
	awsRegion        = "us-east-1"
	s3Bucket         = "myerrorbucket"
	uploadFolderPath = "upload/errors"
)

func UploadError(c *fiber.Ctx) error {
	// Retrieve text input from the form data
	text := c.FormValue("text")
	fmt.Println("text : ", text)

	// Handle text data as needed (e.g., save to a database)

	// Handle image files
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	files := form.File["images"]

	for _, file := range files {
		fileName := file.Filename
		filePath := fmt.Sprintf("%s/%s", uploadFolderPath, fileName)
		err := helpers.UploadToS3(file, filePath, awsRegion, s3Bucket)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to upload file %s to S3", fileName),
			})
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Text and files uploaded successfully",
	})
}
