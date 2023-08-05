package cmd

import (
	"fmt"
	"github.com/purvisb179/profound-crow/cmd/cli"
	"github.com/purvisb179/profound-crow/cmd/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

func Execute() {
	cobra.OnInitialize(initConfig)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err == nil {
		log.Printf("using config.json file found on disk")
	}
}

var rootCmd = &cobra.Command{
	Use:   "pc",
	Short: "PC (Profound Crow) is a CLI for interacting with the Profound Crow API and managing the Asynq worker.",
	Long:  `PC (Profound Crow) provides a set of commands to interact with the Profound Crow API, including calling endpoints and managing tasks. It also allows you to manage the Asynq worker, including starting the worker and handling tasks. This CLI is intended for users of the Profound Crow service.`,
}

func init() {
	rootCmd.AddCommand(
		server.GetStartCmd(),
		cli.GetCheckCmd(),
		cli.GetUploadCmd(),
		cli.GetAuthCmd(),
	)
}
