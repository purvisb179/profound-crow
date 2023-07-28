package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url      string
	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "Send a GET request and print the decoded base64 response",
		Long:  `Send a GET request to the specified URL, decode the base64 response and pretty print it.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			url, _ = cmd.Flags().GetString("url")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
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

func init() {
	checkCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "URL to send the GET request to (required)")
	checkCmd.MarkPersistentFlagRequired("url")
}

func GetCheckCmd() *cobra.Command {
	return checkCmd
}
