package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	queueName string

	clearQueueCmd = &cobra.Command{
		Use:   "clear-queue",
		Short: "Clear a specific queue on the server",
		Long:  `Clear all tasks in a specific queue on the server`,
		PreRun: func(cmd *cobra.Command, args []string) {
			queueName, _ = cmd.Flags().GetString("queue")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			service := "my_cli_app"
			endpoint, err := keyring.Get(service, "endpoint")
			if err != nil {
				return err
			}

			url := fmt.Sprintf("http://%s/clear-queue?queue=%s", endpoint, queueName)

			request, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				return err
			}

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				return err
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Printf("failed to close reader")
				}
			}(response.Body)

			if response.StatusCode == http.StatusOK {
				responseBody, readErr := ioutil.ReadAll(response.Body)
				if readErr != nil {
					return readErr
				}
				fmt.Println(string(responseBody))
			} else {
				return fmt.Errorf("received non-200 response code: %d", response.StatusCode)
			}

			return nil
		},
	}
)

func init() {
	clearQueueCmd.PersistentFlags().StringVarP(&queueName, "queue", "q", "", "Name of the queue to clear (required)")
	clearQueueCmd.MarkPersistentFlagRequired("queue")
}

func GetClearQueueCmd() *cobra.Command {
	return clearQueueCmd
}
