package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

func (s *Server) GenerateDocument(c *fiber.Ctx) error {
	fmt.Println("Start request!!")

	id := c.FormValue("error_id")

	fmt.Println("Printing error_id : ", id)

	filepath := "/errorDocs/" + id

	geminiApiEndPoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="

	var userActions []structures.UserAction
	var errorData structures.Errors

	err := s.Db.Where("error_id = ? AND useful = 1", id).Find(&userActions).Error
	if err != nil {
		fmt.Println("Error fetching user action data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch user action data",
		})
	}

	err = s.Db.Where("id = ?", id).Find(&errorData).Error
	fmt.Println("Printing Title : ", errorData.Title)
	if err != nil {
		fmt.Println("Error fetching error data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch error data",
		})
	}

	requestText := fmt.Sprintf("Assume a user has just solved an error. Here are the details we have: 1. **Title:** %s 2. **Commands/Texts/Code Snippets:** [", errorData.Title)
	for _, val := range userActions {
		requestText += " " + val.TextContent
	}
	requestText += "] Please generate a blog post styled documentation in the manner of a Medium article. The documentation should include the following sections: "
	requestText += "**Introduction** Start with an engaging introduction that provides context about the error and its impact. This should be a paragraph. "
	requestText += "*Understanding the Error* Provide a detailed explanation of the error, including possible causes and symptoms. This should be a paragraph. "
	requestText += "*Troubleshooting Steps* Outline a series of troubleshooting steps that were taken to resolve the error. Each step should be detailed and include any relevant commands or code snippets. "
	requestText += "Commands and code snippets should be wrapped between triple backticks (```). Bullet points should start with *. Subheadings should be marked with a single asterisk (*). "
	requestText += "Ensure a high amount of emphasis on user-generated inputs (user actions) while making the documentation very elaborate. "
	requestText += "The documentation should be written strictly in the first person perspective to make it evident that the person personally solved this error. "
	requestText += "**Conclusion** Summarize the resolution process and offer any additional tips or best practices to avoid similar errors in the future. This should be a paragraph. "
	requestText += "Formatting guidelines: 1. **Introduction**: Place between double asterisks and on a new line with no leading spaces (e.g., **Introduction**). 2. *Subheading*: Place between single asterisks on a new line with no leading spaces (e.g., *Subheading*). "
	requestText += "3. Paragraph: Regular text. 4. Code/Command: Wrap with triple backticks (```). 5. Bullet Point: Start with an asterisk (*) followed by a space. Example structure: **Introduction** This section provides an engaging introduction to the error and its impact. *Understanding the Error* This section details the error, its possible causes, and symptoms. "
	// requestText += "*Troubleshooting Steps* * Step 1: Identify the Issue ``` aws --version ``` * Step 2: Reinstall AWS CLI ``` sudo pip install awscli --force-reinstall --upgrade ``` "
	requestText += "*Troubleshooting Steps* * Step 1: some text and commands or code if necessary  * Step 2: some text and commands or code if necessary"
	requestText += "**Conclusion** This section summarizes the resolution process and offers additional tips or best practices. Ensure the response follows these guidelines strictly with no extra symbols or empty lines."

	requestData := structures.GeminiAIRequest{
		Contents: []structures.Content{
			{
				Parts: []structures.Part{
					{
						Text: requestText,
					},
				},
			},
		},
	}

	fmt.Println("Printing the whole string : ", requestData.Contents[0].Parts[0].Text)

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	fmt.Println("Length of request body : ", len(requestBody))

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	apiKey := os.Getenv("APIKEY")
	geminiApiEndPoint += apiKey

	fmt.Println("Printing endpoint : ", geminiApiEndPoint)

	resp, err := http.Post(geminiApiEndPoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Printing error : ", err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Printing status code : ", resp.StatusCode)
		return err
	}

	var geminiResponse structures.GeminiAIResponse
	err = json.Unmarshal(body, &geminiResponse)
	if err != nil {
		log.Fatal("Encountered the following error while generating the document : ", err)
		return err
	}

	// fmt.Println("Printing the response before extracting the text : ", geminiResponse)

	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		text := geminiResponse.Candidates[0].Content.Parts[0].Text

		// Remove "---" from the text
		text = strings.ReplaceAll(text, "---", "")

		// UPLOAD TO S3
		err := helpers.UploadTextToS3(text, filepath, awsRegion, s3Bucket)
		if err != nil {
			log.Fatal("Unable to upload the error to S3 bucket")
		}

		if err := updateDocFilePath(s.Db, id, filepath); err != nil {
			log.Fatalf("failed to update doc file path: %v", err)
		}

		log.Println("DocFilePath updated successfully")

		// UPDATE THE DOC PATH IN DB
		// result := s.Db.Exec("UPDATE errors SET time_taken = ? WHERE id = ?", timeElapsed, errorID)

		// if result.Error != nil {
		// 	fmt.Println("Error updating record:", result.Error)
		// 	return result.Error
		// }

		fmt.Println("Here is the response from Gemini Ai : ", text)

		c.Set("Content-Type", "application/json")
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":   "success",
			"response": text,
		})
	}

	return err
}

// func (s *Server) updateSavedFileState

func updateDocFilePath(db *gorm.DB, errorID string, newDocFilePath string) error {
	result := db.Model(&structures.Errors{}).Where("id = ?", errorID).Update("doc_file_path", newDocFilePath)
	return result.Error
}
