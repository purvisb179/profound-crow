package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
	"net/http"
	"strconv"
)

type SetTemperatureCommand struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

var SetTempCmd = &cobra.Command{
	Use:   "set-temp [deviceID] [mode] [temp]",
	Short: "Sets the thermostat mode and temperature",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID := args[0]
		mode := args[1]
		temp := args[2]
		return setTemp(deviceID, mode, temp)
	},
}

func setTemp(deviceID string, mode string, temp string) error {
	service := "my_cli_app"

	accessToken, err := keyring.Get(service, "accessToken")
	if err != nil {
		return err
	}

	projectID, err := keyring.Get(service, "projectID")
	if err != nil {
		return err
	}

	commandName, err := mapModeToCommand(mode)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://smartdevicemanagement.googleapis.com/v1/enterprises/%s/devices/%s:executeCommand", projectID, deviceID)

	tempKey, err := mapModeToTemperatureKey(mode)
	if err != nil {
		return err
	}

	tempFloat, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		return fmt.Errorf("could not parse temperature: %v", err)
	}

	command := &SetTemperatureCommand{
		Command: commandName,
		Params:  map[string]interface{}{tempKey: tempFloat},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(command)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuf)
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	fmt.Printf("Mode set to %s and temperature set to %s°C for device: %s\n", mode, temp, deviceID)
	return nil
}

func mapModeToCommand(mode string) (string, error) {
	switch mode {
	case "HEAT":
		return "sdm.devices.commands.ThermostatTemperatureSetpoint.SetHeat", nil
	case "COOL":
		return "sdm.devices.commands.ThermostatTemperatureSetpoint.SetCool", nil
	default:
		return "", fmt.Errorf("invalid mode: %s. available modes are: HEAT, COOL", mode)
	}
}

func mapModeToTemperatureKey(mode string) (string, error) {
	switch mode {
	case "HEAT":
		return "heatCelsius", nil
	case "COOL":
		return "coolCelsius", nil
	default:
		return "", fmt.Errorf("invalid mode: %s. available modes are: HEAT, COOL", mode)
	}
}
