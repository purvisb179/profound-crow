package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var (
	clientID     string
	clientSecret string
	refreshToken string
	oauthURL     string
	endpoint     string
	initCmd      = &cobra.Command{
		Use:   "init",
		Short: "Initialize the CLI with your OAuth credentials and endpoint URL",
		Long:  "Securely store your OAuth client ID, client secret, refresh token, and the endpoint URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(clientID, clientSecret, refreshToken, oauthURL, endpoint)
		},
	}
)

func init() {
	initCmd.Flags().StringVar(&clientID, "client-id", "", "OAuth2 Client ID (required)")
	initCmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth2 Client Secret (required)")
	initCmd.Flags().StringVar(&refreshToken, "refresh-token", "", "OAuth2 Refresh Token (required)")
	initCmd.Flags().StringVar(&oauthURL, "oauth-url", "https://www.googleapis.com/oauth2/v4/token", "Google OAuth Endpoint URL")
	initCmd.Flags().StringVar(&endpoint, "endpoint", "3.14.100.15", "Endpoint URL")
	initCmd.MarkFlagRequired("client-id")
	initCmd.MarkFlagRequired("client-secret")
	initCmd.MarkFlagRequired("refresh-token")
}

func GetInitCmd() *cobra.Command {
	return initCmd
}

func initialize(clientID string, clientSecret string, refreshToken string, oauthURL string, endpoint string) error {
	service := "my_cli_app"

	err := keyring.Set(service, "clientID", clientID)
	if err != nil {
		return err
	}

	err = keyring.Set(service, "clientSecret", clientSecret)
	if err != nil {
		return err
	}

	err = keyring.Set(service, "refreshToken", refreshToken)
	if err != nil {
		return err
	}

	err = keyring.Set(service, "oauthURL", oauthURL)
	if err != nil {
		return err
	}

	err = keyring.Set(service, "endpoint", endpoint)
	if err != nil {
		return err
	}

	fmt.Println("Credentials and endpoint URL stored.")
	return nil
}
