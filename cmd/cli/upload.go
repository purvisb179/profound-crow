package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	filePath          string
	deviceID          string
	name              string
	vacantTemp        int
	occupiedTemp      int
	rampUpTimeSeconds int

	uploadCmd = &cobra.Command{
		Use:   "upload",
		Short: "Upload a file to the server",
		Long:  `Upload a file to the specified URL`,
		PreRun: func(cmd *cobra.Command, args []string) {
			filePath, _ = cmd.Flags().GetString("file")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			service := "my_cli_app"
			endpoint, err := keyring.Get(service, "endpoint")
			if err != nil {
				return err
			}

			url := fmt.Sprintf("http://%s/new-calendar", endpoint)

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Printf("failed to close file")
				}
			}(file)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("myFile", filepath.Base(filePath))
			if err != nil {
				return err
			}
			_, err = io.Copy(part, file)

			// Create the UploadInput struct from the flags
			uploadInput := pkg.UploadInput{
				DeviceID:          deviceID,
				Name:              name,
				VacantTemp:        vacantTemp,
				OccupiedTemp:      occupiedTemp,
				RampUpTimeSeconds: rampUpTimeSeconds,
			}

			// Marshal the struct to JSON
			configJSON, err := json.Marshal(uploadInput)
			if err != nil {
				return err
			}

			// Attach the JSON data to the request
			configField, err := writer.CreateFormField("configuration")
			if err != nil {
				return err
			}
			_, err = configField.Write(configJSON)
			if err != nil {
				return err
			}

			err = writer.Close()
			if err != nil {
				return err
			}

			request, err := http.NewRequest("POST", url, body)
			if err != nil {
				return err
			}
			request.Header.Add("Content-Type", writer.FormDataContentType())

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
	uploadCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path of the file to upload (required)")
	uploadCmd.PersistentFlags().StringVarP(&deviceID, "deviceID", "d", "", "Device ID (required)")
	uploadCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "Name (required)")
	uploadCmd.PersistentFlags().IntVarP(&vacantTemp, "vacantTemp", "v", 0, "Vacant Temperature (required)")
	uploadCmd.PersistentFlags().IntVarP(&occupiedTemp, "occupiedTemp", "o", 0, "Occupied Temperature (required)")
	uploadCmd.PersistentFlags().IntVarP(&rampUpTimeSeconds, "rampUpTime", "r", 0, "Ramp Up Time in Seconds (required)")

	uploadCmd.MarkPersistentFlagRequired("file")
	uploadCmd.MarkPersistentFlagRequired("deviceID")
	uploadCmd.MarkPersistentFlagRequired("name")
	uploadCmd.MarkPersistentFlagRequired("vacantTemp")
	uploadCmd.MarkPersistentFlagRequired("occupiedTemp")
	uploadCmd.MarkPersistentFlagRequired("rampUpTime")
}

func GetUploadCmd() *cobra.Command {
	return uploadCmd
}
