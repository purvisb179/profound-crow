package cmd

import (
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Profound Crow server",
	Long:  `Start the Profound Crow server. This will run the API server, set up the Asynq worker, and get everything ready to accept requests.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func GetServerCmd() *cobra.Command {
	return serverCmd
}
