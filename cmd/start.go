package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Profound Crow server",
	Long:  `Start the Profound Crow server. This will run the API server, set up the Asynq worker, and get everything ready to accept requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		redisAddr := viper.GetString("redisAddr")
		if redisAddr == "" {
			return fmt.Errorf("redis address not specified in configuration file")
		}

		// Remaining server startup code here
		log.Printf("hello world!")
		return nil
	},
}

func GetStartCmd() *cobra.Command {
	return startCmd
}
