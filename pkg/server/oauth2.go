package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random" // Change this to a random string for security
	clientId          = os.Getenv("CLIENT_ID")
	clientSecret      = os.Getenv("CLIENT_SECRET")
)

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  "https://jm4e775kx3.execute-api.us-east-1.amazonaws.com/prod/auth/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (s *Server) HandleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (s *Server) HandleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != oauthStateString {
		return c.Status(http.StatusUnauthorized).SendString("Invalid oauth state")
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	userInfo, err := getUserInfo(token)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(userInfo)
}

func getUserInfo(token *oauth2.Token) (map[string]interface{}, error) {
	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
