package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/pkg"
	"log"
	"time"
)

type AsynqService struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewAsynqService(client *asynq.Client, inspector *asynq.Inspector) *AsynqService {
	return &AsynqService{
		Client:    client,
		Inspector: inspector,
	}
}

func (s *AsynqService) EnqueueTask(task *asynq.Task, durationUntilProcessing time.Duration) error {
	_, err := s.Client.Enqueue(task, asynq.ProcessIn(durationUntilProcessing))
	return err
}

func (s *AsynqService) ListScheduledTasks() ([]*asynq.TaskInfo, error) {
	return s.Inspector.ListScheduledTasks("default", 0, -1)
}

func (s *AsynqService) ProcessAndEnqueueCalendarEvent(payload pkg.CalendarEventPayload) error {
	// Convert payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("could not marshal payload: %v", err)
	}

	// Create the start task
	startTask := asynq.NewTask("CalendarEventPayload", payloadBytes)
	if err != nil {
		return fmt.Errorf("could not create start task: %v", err)
	}

	log.Printf("Start task created successfully. Summary: %s, Process Time: %v", payload.EventSummary, payload.EventStart)

	// Compute the time until processing the start task, but subtract the ramp up time
	durationUntilStart := payload.EventStart.Sub(time.Now()) - time.Second*time.Duration(payload.Configuration.RampUpTimeSeconds)
	if durationUntilStart < 0 {
		return fmt.Errorf("event start in the past or ramp-up time exceeds event start, skipping")
	}

	if err := s.EnqueueTask(startTask, durationUntilStart); err != nil {
		return fmt.Errorf("could not enqueue start task: %v", err)
	}

	log.Printf("Start task enqueued successfully. Duration until processing: %v", durationUntilStart)

	// Create the end task
	endTask := asynq.NewTask("CalendarEventPayload", payloadBytes)
	if err != nil {
		return fmt.Errorf("could not create end task: %v", err)
	}

	log.Printf("End task created successfully. Summary: %s, Process Time: %v", payload.EventSummary, payload.EventEnd)

	// Compute the time until processing the end task
	durationUntilEnd := payload.EventEnd.Sub(time.Now())
	if durationUntilEnd < 0 {
		return fmt.Errorf("event end in the past, skipping")
	}

	if err := s.EnqueueTask(endTask, durationUntilEnd); err != nil {
		return fmt.Errorf("could not enqueue end task: %v", err)
	}

	log.Printf("End task enqueued successfully. Duration until processing: %v", durationUntilEnd)

	return nil
}
