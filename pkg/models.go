package pkg

import "time"

type CalendarEventPayload struct {
	FilePath      string
	EventSummary  string
	EventStart    time.Time
	Configuration map[string]interface{} `json:"configuration"`
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

type UploadInput struct {
	DeviceID          string `json:"device_id"`
	Name              string `json:"name"`
	VacantTemp        int    `json:"vacant_temp"`
	OccupiedTemp      int    `json:"occupied_temp"`
	RampUpTimeSeconds int    `json:"ramp_up_time_seconds"`
}
