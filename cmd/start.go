package cmd

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Profound Crow server",
	Long:  `Start the Profound Crow server. This will run the API server, set up the Asynq worker, and get everything ready to accept requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		redisAddr := viper.GetString("redisAddr")
		if redisAddr == "" {
			return fmt.Errorf("redis address not specified in configuration")
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

		go func() {
			router := chi.NewRouter()
			api.BindRoutes(router)
			if err := http.ListenAndServe(":8888", router); err != nil {
				log.Fatalf("Failed to start API server: %v", err)
			}
		}()

		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}

		return nil
	},
}

func GetStartCmd() *cobra.Command {
	return startCmd
}
