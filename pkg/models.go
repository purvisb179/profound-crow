package pkg

import "time"

type CalendarEventPayload struct { //todo also rename this please. this is supposed to be some higher level meta data of items in the calendar
	FilePath      string
	EventSummary  string
	EventStart    time.Time
	EventEnd      time.Time
	Configuration UploadInput //todo change this naming cause its confusing
}

type CalendarTaskPayload struct {
	DeviceID string
	Temp     int
}

type CalendarTaskCheckResponse struct {
	Payload   CalendarTaskPayload `json:"payload"`
	StartTime time.Time           `json:"startTime"`
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
