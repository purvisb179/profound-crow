package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"net/http"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

var (
	clientID     string
	clientSecret string
	refreshToken string
	authCmd      = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Google OAuth",
		Long:  "Use a refresh token to get a new access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authenticate(clientID, clientSecret, refreshToken)
		},
	}
)

func init() {
	authCmd.Flags().StringVar(&clientID, "client-id", "", "OAuth2 Client ID")
	authCmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth2 Client Secret")
	authCmd.Flags().StringVar(&refreshToken, "refresh-token", "", "OAuth2 Refresh Token")
}

func GetAuthCmd() *cobra.Command {
	return authCmd
}

func authenticate(clientID string, clientSecret string, refreshToken string) error {
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return errors.New("client-id, client-secret and refresh-token are required")
	}

	data := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://www.googleapis.com/oauth2/v4/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %v", string(body))
	}

	var tokenResponse TokenResponse

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Expires in: %s\n", time.Duration(tokenResponse.ExpiresIn)*time.Second)

	service := "my_cli_app"
	user := "user"
	err = keyring.Set(service, user, tokenResponse.AccessToken)
	if err != nil {
		return err
	}

	fmt.Println("Access token saved to keyring.")
	return nil
}
