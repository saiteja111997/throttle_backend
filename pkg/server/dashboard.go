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

	fmt.Println("Printing userId : ", userId)

	var dashboardData []structures.DashboardData

	// Perform the query
	err := s.Db.Raw("SELECT * FROM errors WHERE status = '1' AND status = '2' ORDER BY created_at DESC").Scan(&dashboardData).Error

	if err != nil {
		fmt.Println("Error while fetching from the database: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch from the database",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": dashboardData,
	})

}

func (s *Server) GetDashboardDoc(c *fiber.Ctx) error {

	docFilePath := c.FormValue("doc_file_path")

	fmt.Println("Printing error id : ", docFilePath)

	// Downloading the text from s3

	docContent, err := helpers.DownloadFromS3(docFilePath, awsRegion, s3Bucket)

	if err != nil {
		fmt.Println("Error while fetching from the S3: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch from the S3",
		})
	}

	fmt.Println("Printing the doc content : ", docContent)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": docContent,
	})

}

func (s *Server) PublishDoc(c *fiber.Ctx) error {

	textContent := c.FormValue("content")
	id := c.FormValue("id")

	fmt.Println("Printing content : ", textContent)

	filepath := "/errorDocs/" + id

	err := helpers.UploadTextToS3(textContent, filepath, awsRegion, s3Bucket)
	if err != nil {
		log.Fatal("Unable to upload the error to S3 bucket")
	}

	err = helpers.UpdateDocStatus(s.Db, id)
	if err != nil {
		log.Fatal("Unable to update doc status!!")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": "Successfully uploaded the doc",
	})

}
