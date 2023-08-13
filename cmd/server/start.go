package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/api"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var (
	clientID     string
	clientSecret string
	refreshToken string
	oauthURL     string
	endpoint     string
	projectID    string
	startCmd     = &cobra.Command{
		Use:   "start",
		Short: "Start the Profound Crow server",
		Long:  `Start the Profound Crow server. This will run the API server, set up the Asynq worker, and get everything ready to accept requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !validateToken(clientID, clientSecret, refreshToken) {
				return fmt.Errorf("invalid authorization token")
			}

			redisAddr := viper.GetString("redisAddr")
			if redisAddr == "" {
				return fmt.Errorf("redis address not specified in configuration")
			}

			serverPort := viper.GetString("serverPort")
			if serverPort == "" {
				return fmt.Errorf("server port not specified in configuration")
			}

			srv := asynq.NewServer(
				asynq.RedisClientOpt{Addr: redisAddr},
				asynq.Config{
					Concurrency: 10,
					Queues: map[string]int{
						"critical": 6,
						"default":  3,
						"low":      1,
					},
				},
			)

			mux := asynq.NewServeMux()
			mux.HandleFunc("calendar_event", tasks.HandleCalendarEvent)

			go func() {
				client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
				inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
				apiHandler := api.NewHandler(client, inspector)
				router := chi.NewRouter()
				api.BindRoutes(router, apiHandler)
				if err := http.ListenAndServe(":"+serverPort, router); err != nil {
					log.Fatalf("Failed to start API server: %v", err)
				}
			}()

			if err := srv.Run(mux); err != nil {
				log.Fatalf("could not run server: %v", err)
			}

			return nil
		},
	}
)

func init() {
	startCmd.Flags().StringVar(&clientID, "client-id", "", "OAuth2 Client ID (required)")
	startCmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth2 Client Secret (required)")
	startCmd.Flags().StringVar(&refreshToken, "authorization", "", "OAuth2 Authorization Token (required)")
	startCmd.Flags().StringVar(&oauthURL, "oauth-url", "https://www.googleapis.com/oauth2/v4/token", "Google OAuth Endpoint URL")
	startCmd.Flags().StringVar(&endpoint, "endpoint", "3.14.100.15", "Endpoint URL")
	startCmd.Flags().StringVar(&projectID, "project-id", "5fdd9b9e-c155-4c40-953b-69b576286a62", "Google Cloud Project ID")

	startCmd.MarkFlagRequired("client-id")
	startCmd.MarkFlagRequired("client-secret")
	startCmd.MarkFlagRequired("authorization")
}

func GetStartCmd() *cobra.Command {
	return startCmd
}

func validateToken(clientID string, clientSecret string, refreshToken string) bool {
	url := "https://www.googleapis.com/oauth2/v4/token"
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
