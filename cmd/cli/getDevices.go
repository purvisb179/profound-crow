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

var GetDevicesCmd = &cobra.Command{
	Use:   "get-devices",
	Short: "Fetches a list of devices linked to your Device Access project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return getDevices()
	},
}

func getDevices() error {
	service := "my_cli_app"

	accessToken, err := keyring.Get(service, "accessToken")
	if err != nil {
		return err
	}

	projectID, err := keyring.Get(service, "projectID")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices", projectID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var devices pkg.DeviceResponse
	err = json.Unmarshal(body, &devices)
	if err != nil {
		return err
	}

	for _, device := range devices.Devices {
		fmt.Println("Device Name:", device.Name)
		fmt.Println("Type:", device.Type)
		fmt.Println("Traits:")
		for traitName, traitValue := range device.Traits {
			traitValueBytes, _ := json.MarshalIndent(traitValue, "", "\t")
			fmt.Printf("%s:\n%s\n", traitName, string(traitValueBytes))
		}
		fmt.Println()
	}

	return nil
}
