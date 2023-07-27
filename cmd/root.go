package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "pc",
	Short: "PC (Profound Crow) is a CLI for interacting with the Profound Crow API and managing the Asynq worker.",
	Long:  `PC (Profound Crow) provides a set of commands to interact with the Profound Crow API, including calling endpoints and managing tasks. It also allows you to manage the Asynq worker, including starting the worker and handling tasks. This CLI is intended for users of the Profound Crow service.`,
}

func init() {
	rootCmd.AddCommand()
}
