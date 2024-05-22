package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	helpers "github.com/saiteja111997/throttle_backend/pkg/helper"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

func (s *Server) GenerateDocument(c *fiber.Ctx) error {

	fmt.Println("Start request!!")

	errorID := c.FormValue("error_id")

	fmt.Println("Printing error_id : ", errorID)

	filepath := "/errorDocs/" + errorID

	geminiApiEndPoint := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key="

	var userActions []structures.UserAction
	var errorData structures.Errors

	err := s.Db.Where("error_id = ?", errorID).Find(&userActions).Error

	if err != nil {
		fmt.Println("Error fetching user action data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch user action data",
		})
	}

	err = s.Db.Where("id = ?", errorID).Find(&errorData).Error

	fmt.Println("Printing Title : ", errorData.Title)

	if err != nil {
		fmt.Println("Error fetching error data from db : ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch error data",
		})
	}

	requestData := structures.GeminiAIRequest{
		Contents: []structures.Content{
			{
				Parts: []structures.Part{
					{
						Text: "Assume user has just solved an error. Here are details we have, error description, list of commands, texts and code snippets. You have to generate error documentation based on this information and make sure the documentation should sound as if the user has created the documentation. Here is what an example request looks like title: aws cli version 2 install -'command not found', code/texts/command:[Don't know why AWS bundler could not do it., sudo chmod -R 755 /usr/local/aws-cli].All the sub headings in response should be wrapped between **sub heading**. Please maintain this precise response structure, No Extra # or *. Keep in mind to put a lot of emphasis on user provided code/texts/commands, and provide an elaborate documentation which also covers the concepts related to the error. Now generate the response for this, title: " + errorData.Title + ", code/texts/command:[",
					},
				},
			},
		},
	}

	for _, val := range userActions {
		requestData.Contents[0].Parts[0].Text += " " + val.TextContent
	}

	requestData.Contents[0].Parts[0].Text += "], Please maintain the structure of the document exactly like this one meaning subheadings between **, code between ```, And bullet points starts with * => **Title: -bash: aws: command not found** **Sub Heading: Possible Causes** * AWS CLI is not installed or not in the PATH environment variable. **Sub Heading: Troubleshooting** **1. Install AWS CLI** ```sudo pip install awscli --force-reinstall --upgrade ``` **2. Update PATH Environment Variable** * Open `.bashrc` or `.zshrc` file. * Add the following line: ```export PATH=/usr/local/bin:$PATH ``` * Save and close the file. * Source the file: ```source ~/.bashrc ``` **Sub Heading: Additional Tips** If you encounter issues during installation, try the following: * Uninstall any existing AWS CLI installation: ``` sudo pip uninstall awscli ``` * Ensure you have Python 3.6 or later installed. * Check the installation directory: ```aws --version ``` * If the command still not found, try: ```sudo chmod -R 755 /usr/local/aws-cli ```"

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

	fmt.Println("Printing the response before extracting the text : ", geminiResponse)

	// Extracting text from the response
	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {

		text := geminiResponse.Candidates[0].Content.Parts[0].Text

		// upload the error documentation to s3
		err := helpers.UploadTextToS3(text, filepath, awsRegion, s3Bucket)

		if err != nil {
			log.Fatal("Unable to upload the error to S3 bucket")
		}

		fmt.Println("Here is the response from Gemini Ai : ", geminiResponse.Candidates[0].Content.Parts[0].Text)

		// Now, set the Content-Type header to JSON
		c.Set("Content-Type", "application/json")

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":   "success",
			"response": geminiResponse.Candidates[0].Content.Parts[0].Text,
		})
	}

	return err
}
