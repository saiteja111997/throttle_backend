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

	errorInput := structures.Errors{
		UserID:   userId,
		ID:       uniqueID,
		Title:    title,
		FilePath: uploadFolderPath,
	}

	count := 1

	for _, file := range files {
		// Get the name of the file
		filename := file.Filename
		fmt.Println("Uploaded filename:", filename)
		// Print any additional file information you need
		fmt.Println("File size:", file.Size)
		fmt.Println("MIME type:", file.Header.Get("Content-Type"))

		switch count {
		case 1:
			errorInput.Image1 = uploadFolderPath + "/" + filename
		case 2:
			errorInput.Image2 = uploadFolderPath + "/" + filename
		case 3:
			errorInput.Image3 = uploadFolderPath + "/" + filename
		case 4:
			errorInput.Image4 = uploadFolderPath + "/" + filename
		case 5:
			break
		}
		count++
	}

	go func(e structures.Errors, errDbUpload error, database *gorm.DB) {
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
		"message":    "Text and files uploaded successfully",
		"session_id": uniqueID,
	})
}

func (s *Server) InsertUserActions(c *fiber.Ctx) error {
	text := c.FormValue("text")
	error_id := c.FormValue("error_id")
	user_id := c.FormValue("user_id")
	user_id_int, err := strconv.Atoi(user_id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to convert user id to integer, %v", err),
		})
	}

	e := structures.UserActions{
		TextContent: text,
		UserID:      user_id_int,
		ErrorID:     error_id,
	}

	errDbUpload := s.Db.Create(&e).Error

	fmt.Println("Error : ", errDbUpload)

	if errDbUpload != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to insert into the db, %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Highlighted text uploaded successfully",
	})
}

func (s *Server) DeleteUserAction(c *fiber.Ctx) error {
	// user_id := c.FormValue("user_id")
	user_action_id := c.FormValue("user_action_id")

	integer_id, err := strconv.Atoi(user_action_id)

	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"error":  err,
		})
	}

	// Delete the record from the UserActions table based on ID
	result := s.Db.Where("id = ?", integer_id).Delete(&structures.UserActions{})

	if result.Error != nil {
		fmt.Println("Error deleting record:", result.Error)
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted the user action!!",
	})

}
