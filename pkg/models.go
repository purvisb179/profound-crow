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
