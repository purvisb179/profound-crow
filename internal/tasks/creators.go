package tasks

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/pkg"
	"time"
)

func CreateCalendarEvent(filePath string, eventSummary string, eventStart time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(pkg.CalendarEventPayload{
		FilePath:     filePath,
		EventSummary: eventSummary,
		EventStart:   eventStart,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask("CalendarEventPayload", payload), nil
}
