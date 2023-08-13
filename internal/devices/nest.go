package devices

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type NestService struct {
	OAuthURL string
}

func NewNestService(oauthURL string) *NestService {
	return &NestService{
		OAuthURL: oauthURL,
	}
}

func (s *NestService) ValidateToken(clientID string, clientSecret string, refreshToken string) bool {
	url := s.OAuthURL
	payload := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return false
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Failed to decode token validation response: %v", err)
		return false
	}

	_, exists := result["access_token"]
	return exists
}
