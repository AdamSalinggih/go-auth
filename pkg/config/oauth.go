package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleOauthConfig() *oauth2.Config {
	ClientId := os.Getenv("GOOGLE_CLIENT_ID")
	ClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	RedirectURL := "http://localhost:8080/api/v1/auth/google/callback"
	Scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}

	return &oauth2.Config{
		ClientID:     ClientId,
		ClientSecret: ClientSecret,
		RedirectURL:  RedirectURL,
		Scopes:       Scopes,
		Endpoint:     google.Endpoint,
	}
}
