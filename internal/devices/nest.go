package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"log"
	"net/http"
	"strconv"
)

type NestService struct {
	OAuthURL     string
	ClientID     string
	ClientSecret string
	RefreshToken string
	ProjectID    string
}

func NewNestService(oauthURL, clientID, clientSecret, refreshToken string, projectId string) *NestService {
	return &NestService{
		OAuthURL:     oauthURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		ProjectID:    projectId,
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

func (s *NestService) SetTemperature(deviceID, mode, temp string) error {
	commandName, err := mapModeToCommand(mode)
	if err != nil {
		return err
	}

	tempKey, err := mapModeToTemperatureKey(mode)
	if err != nil {
		return err
	}

	tempFloat, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return fmt.Errorf("could not parse temperature: %v", err)
	}

	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices/%s:executeCommand", s.ProjectID, deviceID)

	command := &pkg.SetTemperatureCommand{
		Command: commandName,
		Params:  map[string]interface{}{tempKey: tempFloat},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(command)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
	if err != nil {
		return err
	}

	token, err := s.getRefreshedToken() //todo should probably do this somewhere else above? like at caller level?
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	log.Printf("Mode set to %s and temperature set to %s°C for device: %s\n", mode, temp, deviceID)
	return nil
}

func mapModeToCommand(mode string) (string, error) {
	switch mode {
	case "HEAT":
		return "sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat", nil
	case "COOL":
		return "sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool", nil
	default:
		return "", fmt.Errorf("invalid mode: %s. available modes are: HEAT, COOL", mode)
	}
}

func mapModeToTemperatureKey(mode string) (string, error) {
	switch mode {
	case "HEAT":
		return "heatCelsius", nil
	case "COOL":
		return "coolCelsius", nil
	default:
		return "", fmt.Errorf("invalid mode: %s. available modes are: HEAT, COOL", mode)
	}
}
