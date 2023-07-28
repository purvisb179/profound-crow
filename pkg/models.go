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
