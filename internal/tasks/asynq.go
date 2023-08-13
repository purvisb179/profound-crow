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
	// Starting Task payload
	startPayload := CalendarTaskPayload{
		DeviceID: payload.Configuration.DeviceID,
		Temp:     payload.Configuration.OccupiedTemp,
	}
	startPayloadBytes, err := json.Marshal(startPayload)
	if err != nil {
		return fmt.Errorf("could not marshal start payload: %v", err)
	}

	// Calculate duration until processing for the starting task
	startDuration := payload.EventStart.Sub(time.Now()) - time.Second*time.Duration(payload.Configuration.RampUpTimeSeconds)
	if startDuration < 0 {
		return fmt.Errorf("event in the past or ramp-up time exceeds event start, skipping starting task")
	}

	// Create and enqueue the starting task
	startTask := asynq.NewTask("calendar_event", startPayloadBytes)
	if err := s.EnqueueTask(startTask, startDuration); err != nil {
		return fmt.Errorf("could not enqueue starting task: %v", err)
	}
	log.Printf("Starting task enqueued successfully. Duration until processing: %v", startDuration)

	// Ending Task payload
	endPayload := CalendarTaskPayload{
		DeviceID: payload.Configuration.DeviceID,
		Temp:     payload.Configuration.VacantTemp,
	}
	endPayloadBytes, err := json.Marshal(endPayload)
	if err != nil {
		return fmt.Errorf("could not marshal end payload: %v", err)
	}

	// Calculate duration until processing for the ending task
	endDuration := payload.EventEnd.Sub(time.Now())
	if endDuration < 0 {
		return fmt.Errorf("event end time is in the past, skipping ending task")
	}

	// Create and enqueue the ending task
	endTask := asynq.NewTask("end_calendar_event", endPayloadBytes) // Name change for clarity
	if err := s.EnqueueTask(endTask, endDuration); err != nil {
		return fmt.Errorf("could not enqueue ending task: %v", err)
	}
	log.Printf("Ending task enqueued successfully. Duration until processing: %v", endDuration)

	return nil
}

// ClearQueue deletes all tasks in the specified queue.
func (s *AsynqService) ClearQueue(queueName string) error {
	deleted, err := s.Inspector.DeleteAllScheduledTasks(queueName)
	if err != nil {
		return fmt.Errorf("could not clear queue %s: %v", queueName, err)
	}
	log.Printf("Successfully deleted %d tasks from queue %s", deleted, queueName)
	return nil
}
