package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io"
	"io/ioutil"
	"log"
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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Printf("failed to close reader")
				}
			}(response.Body)

			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("received non-200 response code: %d", response.StatusCode)
			}

			body, readErr := ioutil.ReadAll(response.Body)
			if readErr != nil {
				return readErr
			}

			var events []pkg.Event
			unmarshalErr := json.Unmarshal(body, &events)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			for _, event := range events {
				decodedBytes, decodeErr := base64.StdEncoding.DecodeString(event.Payload)
				if decodeErr != nil {
					return decodeErr
				}

				var payload pkg.CalendarEventPayload
				unmarshalPayloadErr := json.Unmarshal(decodedBytes, &payload)
				if unmarshalPayloadErr != nil {
					return unmarshalPayloadErr
				}

				prettyJSON, err := json.MarshalIndent(payload, "", "\t")
				if err != nil {
					return err
				}

				fmt.Println(string(prettyJSON))
			}
			return nil
		},
	}
)

func GetCheckCmd() *cobra.Command {
	return checkCmd
}
