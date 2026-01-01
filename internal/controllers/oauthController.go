package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/adamhaiqal/go-auth/internal/config"
	"github.com/adamhaiqal/go-auth/internal/models"
	"github.com/adamhaiqal/go-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// GoogleOAuthConfig holds OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
}

// GetGoogleOAuthConfig returns the Google OAuth configuration
func GetGoogleOAuthConfig() *GoogleOAuthConfig {
	return &GoogleOAuthConfig{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"), // e.g., "http://localhost:8080/api/v1/oauth/google/callback"
		AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
		UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

// GoogleLogin initiates the OAuth flow by redirecting to Google
func GoogleLogin(c *gin.Context) {
	config := GetGoogleOAuthConfig()
	if config.ClientID == "" || config.ClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth configuration missing"})
		return
	}

	// Generate state token for CSRF protection (you might want to store this in session/redis)
	state := generateStateToken()
	// In production, store this state in session/redis and validate it in callback

	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid%%20email%%20profile&state=%s&access_type=offline&prompt=consent",
		config.AuthURL,
		config.ClientID,
		config.RedirectURL,
		state,
	)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// GoogleCallback handles the callback from Google OAuth
func GoogleCallback(c *gin.Context) {
	config := GetGoogleOAuthConfig()

	// Get authorization code from query parameter
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Validate state token (in production, check against stored state)
	_ = state // TODO: Validate state token

	// Exchange authorization code for access token
	token, err := exchangeCodeForToken(config, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token", "details": err.Error()})
		return
	}

	// Get user info from Google
	userInfo, err := getUserInfoFromGoogle(config, token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info", "details": err.Error()})
		return
	}

	// Find or create user in database
	account, err := findOrCreateOAuthUser(userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user", "details": err.Error()})
		return
	}

	// Generate JWT token (same as regular signin)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      account.ID,
		"username": account.Username,
		"email":    account.Email,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := jwtToken.SignedString(utils.GetJWTKey())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set cookie (same as regular signin)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenString, 3600, "", "", false, true)

	// Redirect to home or return success
	c.JSON(http.StatusOK, gin.H{
		"message":  "OAuth signin successful",
		"username": account.Username,
		"email":    account.Email,
	})
}

// TokenResponse represents the token response from Google
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

// GoogleUserInfo represents user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// exchangeCodeForToken exchanges the authorization code for an access token
func exchangeCodeForToken(config *GoogleOAuthConfig, code string) (*TokenResponse, error) {
	req, err := http.NewRequest("POST", config.TokenURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("client_id", config.ClientID)
	q.Add("client_secret", config.ClientSecret)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", config.RedirectURL)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// getUserInfoFromGoogle fetches user information from Google
func getUserInfoFromGoogle(config *GoogleOAuthConfig, accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequest("GET", config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// findOrCreateOAuthUser finds an existing user or creates a new one from OAuth info
func findOrCreateOAuthUser(userInfo *GoogleUserInfo) (*models.Account, error) {
	var account models.Account

	// Try to find user by email
	err := config.DB.Where("email = ?", userInfo.Email).First(&account).Error

	if err != nil {
		// User doesn't exist, create new account
		// Generate a random password (OAuth users won't use it, but it's required by your model)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(fmt.Sprintf("oauth_%s_%d", userInfo.ID, time.Now().Unix())), bcrypt.DefaultCost)

		account = models.Account{
			Username:   userInfo.Email, // Use email as username, or generate one
			Password:   string(hashedPassword),
			Email:      userInfo.Email,
			FirstName:  userInfo.GivenName,
			LastName:   userInfo.FamilyName,
			IsVerified: userInfo.VerifiedEmail,
			// Set default values for required fields
			Address:     "N/A",
			HomePhone:   "N/A",
			MobilePhone: "N/A",
			WorkPhone:   "N/A",
			StateCode:   "N/A",
			ZipCode:     "N/A",
			Country:     "N/A",
		}

		if err := config.DB.Create(&account).Error; err != nil {
			return nil, err
		}
	} else {
		// User exists, update verification status if needed
		if userInfo.VerifiedEmail && !account.IsVerified {
			account.IsVerified = true
			config.DB.Save(&account)
		}
	}

	return &account, nil
}

// generateStateToken generates a random state token for CSRF protection
func generateStateToken() string {
	// In production, use crypto/rand for secure random generation
	return fmt.Sprintf("state_%d", time.Now().UnixNano())
}
