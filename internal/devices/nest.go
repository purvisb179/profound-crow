package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NestService struct {
	OAuthURL     string
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func NewNestService(oauthURL, clientID, clientSecret, refreshToken string) *NestService {
	return &NestService{
		OAuthURL:     oauthURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
	}
}

// Obtain refreshed token
func (s *NestService) getRefreshedToken() (string, error) {
	url := s.OAuthURL
	payload := map[string]string{
		"client_id":     s.ClientID,
		"client_secret": s.ClientSecret,
		"refresh_token": s.RefreshToken,
		"grant_type":    "refresh_token",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	accessToken, exists := result["access_token"]
	if !exists {
		return "", fmt.Errorf("no access_token in response")
	}

	return accessToken.(string), nil
}

func (s *NestService) ValidateToken() bool {
	_, err := s.getRefreshedToken()
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return false
	}
	return true
}
