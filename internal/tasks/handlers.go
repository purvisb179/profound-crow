package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/devices"
	"github.com/purvisb179/profound-crow/pkg"
	"log"
	"time"
)

type TaskHandler struct {
	NestService *devices.NestService
}

func (th *TaskHandler) HandleCalendarEvent(ctx context.Context, t *asynq.Task) error {
	var p pkg.CalendarEventPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Use th.NestService here

	log.Println("Event Summary:", p.EventSummary)
	log.Println("Event Start:", p.EventStart.Format(time.RFC3339))

	return nil
}
