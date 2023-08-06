package pkg

import "time"

type CalendarEventPayload struct {
	FilePath     string
	EventSummary string
	EventStart   time.Time
}

type Event struct {
	ID      string
	Payload string
	Type    string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type Device struct {
	Name            string
	Type            string
	Traits          map[string]interface{}
	ParentRelations []map[string]string
}

type DeviceResponse struct {
	Devices []Device
}

type SetTemperatureCommand struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}
