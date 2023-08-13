package cli

import (
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"net/http"
)

var (
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Send a GET request and print the decoded base64 response",
		Long:  "Send a GET request to the endpoint, decode the base64 response and pretty print it.",
		RunE: func(cmd *cobra.Command, args []string) error {
			service := "my_cli_app"
			endpoint, err := keyring.Get(service, "endpoint")
			if err != nil {
				return err
			}

			url := fmt.Sprintf("http://%s/check-queue", endpoint)

			response, err := http.Get(url)
			if err != nil {
				return err
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("received non-200 response code: %d", response.StatusCode)
			}

			body, readErr := ioutil.ReadAll(response.Body)
			if readErr != nil {
				return readErr
			}

			// Print the raw response if verbose flag is set
			if verbose {
				fmt.Println("Raw Response:", string(body))
			}

			var taskDetails []pkg.CalendarTaskCheckResponse
			unmarshalErr := json.Unmarshal(body, &taskDetails)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			for _, detail := range taskDetails {
				prettyJSON, err := json.MarshalIndent(detail, "", "\t")
				if err != nil {
					return err
				}

				fmt.Println(string(prettyJSON))
			}
			return nil
		},
	}
)

func init() {
	// Add the verbose flag to the checkCmd
	checkCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print raw server response")
}

func GetCheckCmd() *cobra.Command {
	return checkCmd
}
