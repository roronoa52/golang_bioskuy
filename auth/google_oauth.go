package auth

import (
	"bioskuy/exception"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2V2 "google.golang.org/api/oauth2/v2"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "random"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleLoginURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state)
}

func GetGoogleUser(code string, c *gin.Context) (*oauth2.Token, *oauth2V2.Userinfo, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return nil, nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	oauth2Service, err := oauth2V2.New(googleOauthConfig.Client(context.Background(), token))
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return nil, nil, fmt.Errorf("oauth2 service creation failed: %s", err.Error())
	}

	userinfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return nil, nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	return token, userinfo, nil
}
