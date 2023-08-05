package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"net/http"
)

var (
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Google OAuth",
		Long:  "Use a refresh token to get a new access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authenticate()
		},
	}
)

func GetAuthCmd() *cobra.Command {
	return authCmd
}

func authenticate() error {
	service := "my_cli_app"

	clientID, err := keyring.Get(service, "clientID")
	if err != nil {
		return err
	}

	clientSecret, err := keyring.Get(service, "clientSecret")
	if err != nil {
		return err
	}

	refreshToken, err := keyring.Get(service, "refreshToken")
	if err != nil {
		return err
	}

	endpoint, err := keyring.Get(service, "endpoint")
	if err != nil {
		return err
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

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
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

	var tokenResponse pkg.TokenResponse

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return err
	}

	fmt.Printf("Expires in: %s\n")

	user := "user"
	err = keyring.Set(service, user, tokenResponse.AccessToken)
	if err != nil {
		return err
	}

	fmt.Println("Access token saved to keyring.")
	return nil
}
