package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v2"
	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

const (
	awsRegion = "us-east-1"
	s3Bucket  = "myerrorbucket"
)

func (s *Server) GetRawErrorDocs(c *fiber.Ctx) error {
	errorID := c.FormValue("error_id")

	var result structures.RawErrorResponse
	var errorInfo structures.Errors
	var userActions []structures.UserAction
	var errorResponse error

	var wg sync.WaitGroup
	wg.Add(2)

	// fetch data from db using go routines
	go func(errorID string) {
		defer wg.Done()
		err := s.Db.Where("error_id = ?", errorID).Order("created_at asc").Find(&userActions).Error
		if err != nil {
			fmt.Println("Error fetching user action data:", err)
			errorResponse = err
		}
		result.UserActions = userActions
	}(errorID)

	go func(errorID string) {
		defer wg.Done()
		err := s.Db.Where("id = ?", errorID).Find(&errorInfo).Error
		if err != nil {
			fmt.Println("Error fetching error info data:", err)
			errorResponse = err
		}
		result.ErrorInfo = errorInfo
	}(errorID)

	wg.Wait()

	if errorResponse != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch data",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func (s *Server) GetImagesFromS3(c *fiber.Ctx) error {
	errorID := c.FormValue("error_id")
	var errorInfo structures.Errors
	imageArr := []string{}

	err := s.Db.Where("id = ?", errorID).Find(&errorInfo).Error
	if err != nil {
		fmt.Println("Error fetching user action data:", err)
	}

	for i := 1; i <= 4; i++ {
		if i == 1 {
			if errorInfo.Image1 != "" {
				imageArr = append(imageArr, errorInfo.Image1)
			}
		} else if i == 2 {
			if errorInfo.Image1 != "" {
				imageArr = append(imageArr, errorInfo.Image2)
			}
		} else if i == 3 {
			if errorInfo.Image1 != "" {
				imageArr = append(imageArr, errorInfo.Image3)
			}
		} else {
			if errorInfo.Image1 != "" {
				imageArr = append(imageArr, errorInfo.Image4)
			}
		}
	}

	imageCh := make(chan []byte)

	var wg sync.WaitGroup

	for _, val := range imageArr {
		if val != "" {
			fmt.Println("Here is the path to image:", val)

			wg.Add(1)
			go func(imageCh chan []byte, filePath string) {
				defer wg.Done()
				imageData, err := helpers.DownloadFromS3(filePath, awsRegion, s3Bucket)
				if err != nil {
					fmt.Println("Unable to fetch the image data due to error:", err)
				}
				imageCh <- imageData
			}(imageCh, val)
		}
	}

	go func() {
		wg.Wait()
		close(imageCh)
	}()

	result := []map[string]string{}
	for i := 0; i < len(imageArr); i++ {
		// key := "Image" + strconv.Itoa(i+1)

		// Encode []byte to base64 string
		base64String := base64.StdEncoding.EncodeToString(<-imageCh)
		// fmt.Println("Base64 string:", base64String)

		result = append(result, map[string]string{
			"image": base64String,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": result,
	})
}

func (s *Server) GetLatestRawError(c *fiber.Ctx) error {

	errorID := c.FormValue("error_id")

	var result []structures.UserAction
	err := s.Db.Where("error_id = ?", errorID).Order("created_at asc").Find(&result).Error
	if err != nil {
		fmt.Println("Error fetching user action data:", err)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "success",
		"result": result,
	})

}
