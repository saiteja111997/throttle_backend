package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

func (s *Server) GetDashboard(c *fiber.Ctx) error {

	userId := c.FormValue("user_id")
	docType := c.FormValue("status")

	fmt.Println("Printing userId : ", userId)
	fmt.Println("Printing docType : ", docType)

	var dashboardData []structures.DashboardData

	// Perform the query
	err := s.Db.Raw("SELECT * FROM errors WHERE status = ? AND user_id = ? ORDER BY created_at DESC", docType, userId).Scan(&dashboardData).Error

	if err != nil {
		fmt.Println("Error while fetching from the database: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch from the database",
		})
	}

	fmt.Println("Length of dashboard data : ", len(dashboardData))

	for _, val := range dashboardData {
		val.Status = docType
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": dashboardData,
	})

}

func (s *Server) GetDashboardDoc(c *fiber.Ctx) error {

	docFilePath := c.FormValue("doc_file_path")
	error_id := c.FormValue("error_id")

	fmt.Println("Printing doc file path : ", docFilePath)
	fmt.Println("Error id : ", error_id)

	var errorData structures.Errors

	// Downloading the text from s3

	docContent, err := helpers.DownloadFromS3(docFilePath, awsRegion, s3Bucket)

	if err != nil {
		fmt.Println("Error while fetching from the S3: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch from the S3",
		})
	}

	docContentString := string(docContent)

	err = s.Db.Where("id = ?", error_id).Find(&errorData).Error
	fmt.Println("Printing Title : ", errorData.Title)
	if err != nil {
		fmt.Println("Error fetching error data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch error data",
		})
	}

	var userData structures.Users

	err = s.Db.Where("id = ?", errorData.UserID).Find(&userData).Error
	if err != nil {
		fmt.Println("Error fetching error data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch user data",
		})
	}

	var name string

	if userData.Username == "" || len(userData.Username) == 0 {
		name = userData.Email
	} else {
		name = userData.Username
	}

	fmt.Println("Printing name and profile data from db : ", name, userData.ProfilePic)

	// Format the article published date in "Jan 1st, 2024" style
	articlePublishedDate := helpers.FormatDate(errorData.UpdatedAt)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"result":  docContentString,
		"title":   errorData.Title,
		"created": articlePublishedDate,
		"user":    name,
		"picture": userData.ProfilePic,
	})

}

func (s *Server) PublishDoc(c *fiber.Ctx) error {

	textContent := c.FormValue("content")
	id := c.FormValue("id")
	status := c.FormValue("status")

	// fmt.Println("Printing content : ", textContent)
	fmt.Println("Printing status : ", status)
	fmt.Println("Printing error id : ", id)
	filepath := "/errorDocs/" + id

	err := helpers.UploadTextToS3(textContent, filepath, awsRegion, s3Bucket)
	if err != nil {
		log.Fatal("Unable to upload the error to S3 bucket")
	}

	err = helpers.UpdateDocStatus(s.Db, id, status)
	if err != nil {
		log.Fatal("Unable to update doc status!!", err.Error())
	}
	fmt.Println("Successfully uploaded the document to s3!!")

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": "Successfully uploaded the doc",
	})

}

func (s *Server) SaveDoc(c *fiber.Ctx) error {

	textContent := c.FormValue("content")
	id := c.FormValue("error_id")

	// fmt.Println("Printing content : ", textContent)
	fmt.Println("Printing error id : ", id)
	filepath := "/errorDocs/" + id

	err := helpers.UploadTextToS3(textContent, filepath, awsRegion, s3Bucket)
	if err != nil {
		log.Fatal("Unable to upload the error to S3 bucket")
	}
	fmt.Println("Successfully uploaded the document to s3!!")

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": "Successfully uploaded the doc",
	})

}

func (s *Server) DeleteDoc(c *fiber.Ctx) error {
	errorID := c.FormValue("error_id")
	docFilePath := c.FormValue("doc_file_path")

	err := helpers.DeleteFromS3(docFilePath, awsRegion, s3Bucket)
	if err != nil {
		log.Fatal("Unable to delete the error from S3 bucket")
	}

	err = helpers.DeleteDocFromDB(s.Db, errorID)
	if err != nil {
		log.Fatal("Unable to delete doc from db!!", err.Error())
	}
	fmt.Println("Successfully deleted the document from s3 and db!!")

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Successfully deleted the user action!!",
	})
}
