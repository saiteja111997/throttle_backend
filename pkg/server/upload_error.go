package server

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

type Server struct {
	Db *gorm.DB
}

var wg sync.WaitGroup

const (
	awsRegion = "us-east-1"
	s3Bucket  = "myerrorbucket"
)

func (s *Server) UploadError(c *fiber.Ctx) error {

	uniqueID := uuid.New().String()
	fmt.Println("Error id is : ", uniqueID)

	userIdString := c.FormValue("userId")
	userId, err := strconv.Atoi(userIdString)
	title := c.FormValue("text")

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse user id!!",
		})
	}

	// Retrieve text input from the form data
	text := c.FormValue("text")
	fmt.Println("text : ", text)

	// Handle image files
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	files := form.File["images"]

	uploadFolderPath := "errors"

	// Get the last three letters
	lastThreeLetters := uniqueID[len(uniqueID)-3:]

	uploadFolderPath += "/" + lastThreeLetters + "/" + uniqueID

	var errFileUpload, errDbUpload error

	wg.Add(2)

	errorInput := structures.Error{
		UserID: userId,
		ID:     uniqueID,
		Title:  title,
	}

	go func(e structures.Error, errDbUpload error, database *gorm.DB) {
		defer wg.Done()
		response := database.Create(&e)
		errDbUpload = response.Error
	}(errorInput, errDbUpload, s.Db)

	go func(files []*multipart.FileHeader, errFileUpload error) {
		defer wg.Done()
		for _, file := range files {
			fileName := file.Filename
			filePath := fmt.Sprintf("%s/%s", uploadFolderPath, fileName)
			errFileUpload = helpers.UploadToS3(file, filePath, awsRegion, s3Bucket)
		}
	}(files, errFileUpload)

	wg.Wait()

	if errDbUpload != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create entries in the db : %v", errDbUpload),
		})
	} else {
		fmt.Println("Data successfully uploaded to Db!!")
	}

	if errFileUpload != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to upload files to S3 : %v", errFileUpload),
		})
	} else {
		fmt.Println("Files successfully added to s3 bucket")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Text and files uploaded successfully",
	})
}
