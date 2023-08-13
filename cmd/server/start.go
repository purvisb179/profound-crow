package server

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/api"
	"github.com/purvisb179/profound-crow/internal/devices"
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
			nestService := devices.NewNestService(oauthURL, clientID, clientSecret, refreshToken, projectID)
			if !nestService.ValidateToken() {
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
			handlerInstance := &tasks.TaskHandler{NestService: nestService}
			mux.HandleFunc("calendar_event", handlerInstance.HandleCalendarEvent)

			go func() {
				client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
				inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: redisAddr})
				asynqService := tasks.NewAsynqService(client, inspector)
				apiHandler := api.NewHandler(asynqService)
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
	startCmd.Flags().StringVar(&refreshToken, "refresh-token", "", "refresh token (required)")
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
