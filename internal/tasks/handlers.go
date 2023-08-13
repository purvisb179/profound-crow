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

	// Determine the mode based on temperature
	var mode string
	if p.Temp <= 65 {
		mode = "COOL"
	} else {
		mode = "HEAT"
	}

	// Use th.NestService to set temperature for a device
	if err := th.NestService.SetTemperature(p.DeviceID, mode, fmt.Sprintf("%d", p.Temp)); err != nil {
		return fmt.Errorf("failed to set temperature: %v: %w", err, asynq.SkipRetry)
	}

	log.Printf("Device ID: %s, Temp: %d, Mode: %s\n", p.DeviceID, p.Temp, mode)

	return nil
}
