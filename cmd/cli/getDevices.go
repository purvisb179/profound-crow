package cli

import (
	"encoding/json"
	"fmt"
	"github.com/purvisb179/profound-crow/pkg"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	verbose       bool
	GetDevicesCmd = &cobra.Command{
		Use:   "get-devices",
		Short: "Fetches a list of devices linked to your Device Access project",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getDevices(verbose)
		},
	}
)

func init() {
	GetDevicesCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

func getDevices(verbose bool) error {
	service := "my_cli_app"

	accessToken, err := keyring.Get(service, "accessToken")
	if err != nil {
		log.Printf("Error getting access token: %v", err)
		return err
	}

	projectID, err := keyring.Get(service, "projectID")
	if err != nil {
		log.Printf("Error getting project ID: %v", err)
		return err
	}

	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices", projectID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK HTTP status: %v", resp.Status)
		log.Printf("Response body: %v", string(body))
		log.Printf(accessToken)
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var devices pkg.DeviceResponse
	err = json.Unmarshal(body, &devices)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return err
	}

	for _, device := range devices.Devices {
		fmt.Println("Device Name:", device.Name)
		fmt.Println("Type:", device.Type)
		if verbose {
			fmt.Println("Traits:")
			for traitName, traitValue := range device.Traits {
				traitValueBytes, _ := json.MarshalIndent(traitValue, "", "\t")
				fmt.Printf("%s:\n%s\n", traitName, string(traitValueBytes))
			}
			fmt.Println()
		}
	}

	return nil
}
