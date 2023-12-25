package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/saiteja111997/throttle_backend/pkg/structures"
)

const geminiApiEndPoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=AIzaSyALm_rgqVG1DTx2ChONsxf_I5OpNc2cbQs"

func (s *Server) GenerateDocument(c *fiber.Ctx) error {

	fmt.Println("Start request!!")

	errorID := c.FormValue("error_id")
	title := c.FormValue("title")

	var userActions []structures.UserAction

	err := s.Db.Where("error_id = ?", errorID).Find(&userActions).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"error":  "Failed to fetch data",
		})
	}

	requestData := structures.GeminiAIRequest{
		Contents: []structures.Content{
			{
				Parts: []structures.Part{
					{
						Text: "Assume user has just solved an error. Here are details we have, error description, list of commands, texts and code snippets. You have to generate error documentation based on this information. Here is what an example request looks like title: aws cli version 2 install -'command not found' error, code/texts/command:[Don't know why AWS bundler could not do it., sudo chmod -R 755 /usr/local/aws-cli].All the sub headings in response should be wrapped between **sub heading**. Now generate the response for this, title: " + title + ", code/texts/command:[",
					},
				},
			},
		},
	}

	previous_string := requestData.Contents[0].Parts[0].Text

	for _, val := range userActions {
		previous_string += val.TextContent
	}

	previous_string += "]"

	fmt.Println("Printing the whole string : ", previous_string)

	requestData.Contents[0].Parts[0].Text = previous_string

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	resp, err := http.Post(geminiApiEndPoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}

	var geminiResponse structures.GeminiAIResponse
	err = json.Unmarshal(body, &geminiResponse)
	if err != nil {
		return err
	}

	// Extracting text from the response
	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"status":   "success",
			"response": geminiResponse.Candidates[0].Content.Parts[0].Text,
		})
	}

	return err
}
