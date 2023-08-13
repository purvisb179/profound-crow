package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/devices"
	"log"
)

type CalendarTaskPayload struct {
	DeviceID string
	Temp     int
}

type TaskHandler struct {
	NestService *devices.NestService
}

func (th *TaskHandler) HandleCalendarEvent(ctx context.Context, t *asynq.Task) error {
	var p CalendarTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Use th.NestService with the payload data (p.DeviceID and p.Temp)
	// For example, if you had a method in NestService to set temperature for a device:
	// th.NestService.SetTemperature(p.DeviceID, p.Temp)

	log.Printf("Device ID: %s, Temp: %d\n", p.DeviceID, p.Temp)

	return nil
}
