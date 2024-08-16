package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

type Server struct {
	Db *gorm.DB
}

func (s *Server) UploadError(c *fiber.Ctx) error {

	uniqueID := uuid.New().String()
	fmt.Println("Error id is : ", uniqueID)

	userIdString := c.FormValue("userId")

	fmt.Println("Printing userID : ", userIdString)

	userId, err := strconv.Atoi(userIdString)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse user id!!",
		})
	}

	fmt.Println("Printing userID : ", userId)

	// Retrieve text input from the form data
	title := c.FormValue("text")
	fmt.Println("title : ", title)

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

	errorInput := structures.Errors{
		UserID: userId,
		ID:     uniqueID,
		Title:  title,
		Status: 0,
		Type:   0,
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
		}
		count++
	}

	errDbUpload := s.Db.Create(&errorInput).Error

	fmt.Println("Print error : ", errDbUpload)

	if errDbUpload != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create entries in the db : %v", errDbUpload),
		})
	} else {
		fmt.Println("Data successfully uploaded to Db!!")
	}

	for _, file := range files {
		fileName := file.Filename
		filePath := fmt.Sprintf("%s/%s", uploadFolderPath, fileName)
		errFileUpload := helpers.UploadToS3(file, filePath, awsRegion, s3Bucket)

		if errFileUpload != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to upload files to S3 : %v", errFileUpload),
			})
		} else {
			fmt.Println("Files successfully added to s3 bucket")
		}

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

	fmt.Println("Printing user id in user action api : ", user_id)

	user_id_int, err := strconv.Atoi(user_id)

	if err != nil {
		fmt.Println("Failed to convert user id to integer : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to convert user id to integer, %v", err),
		})
	}

	// startContainer := c.FormValue("startContainer")
	// endContainer := c.FormValue("endContainer")
	// startOffset := c.FormValue("startOffset")
	// endOffset := c.FormValue("endOffset")
	currentURL := c.FormValue("currentURL")

	// fmt.Println("startContainer, endContainer, startOffset, endOffset : ", startContainer, endContainer, startOffset, endOffset)
	fmt.Println("currentURL : ", currentURL)

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

func (s *Server) ValidateUserAction(c *fiber.Ctx) error {
	user_action_id := c.FormValue("user_action_id")
	isUseful := c.FormValue("useful")

	fmt.Println("Printing useful : ", isUseful)
	fmt.Println("Printing user_action_id : ", user_action_id)

	if isUseful != "0" && isUseful != "1" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid useful field.",
		})
	}

	integer_id, err := strconv.Atoi(user_action_id)

	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"error":  err,
		})
	}

	result := s.Db.Exec("UPDATE user_actions SET useful = ? WHERE id = ?", isUseful, integer_id)

	if result.Error != nil {
		fmt.Println("Error deleting record:", result.Error)
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully validated the user action!!",
	})
}

func (s *Server) UpdateErrorState(c *fiber.Ctx) error {

	errorID := c.FormValue("error_id")
	timeElapsed := c.FormValue("elapsed_time")

	fmt.Println("Printing errorID : ", errorID)
	fmt.Println("Printing timeElapsed : ", timeElapsed)

	result := s.Db.Exec("UPDATE errors SET time_taken = ? WHERE id = ?", timeElapsed, errorID)

	if result.Error != nil {
		fmt.Println("Error updating record:", result.Error)
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully validated the user action!!",
	})
}

func (s *Server) UpdateFinalState(c *fiber.Ctx) error {

	errorID := c.FormValue("error_id")
	finalState := c.FormValue("status")

	finalStateInt, err := strconv.Atoi(finalState)

	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "failure",
			"error":  err,
		})
	}

	fmt.Println("Printing errorID : ", errorID)
	fmt.Println("Printing timeElapsed : ", finalStateInt)

	result := s.Db.Exec("UPDATE errors SET status = ? WHERE id = ?", finalState, errorID)

	if result.Error != nil {
		fmt.Println("Error updating record:", result.Error)
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully validated the user action!!",
	})
}

func (s *Server) GetUnresolvedJourneys(c *fiber.Ctx) error {

	// Define a slice to hold the query results
	var unresolvedJourneys []structures.GetUnresolvedJourneys

	// Perform the query
	err := s.Db.Raw("SELECT * FROM errors WHERE status = '0' ORDER BY created_at DESC LIMIT 3").Scan(&unresolvedJourneys).Error

	if err != nil {
		fmt.Println("Error while fetching from the database: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch from the database",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": unresolvedJourneys,
	})

}
