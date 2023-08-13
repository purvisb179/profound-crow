package tasks

import (
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
	task, err := CreateCalendarEvent(payload.FilePath, payload.EventSummary, payload.EventStart, payload.Configuration)
	if err != nil {
		return fmt.Errorf("could not create task: %v", err)
	}

	log.Printf("Task created successfully. Summary: %s, Process Time: %v", payload.EventSummary, payload.EventStart)

	durationUntilProcessing := payload.EventStart.Sub(time.Now())
	if durationUntilProcessing < 0 {
		return fmt.Errorf("event in the past, skipping")
	}

	if err := s.EnqueueTask(task, durationUntilProcessing); err != nil {
		return fmt.Errorf("could not enqueue task: %v", err)
	}

	log.Printf("Task enqueued successfully. Duration until processing: %v", durationUntilProcessing)

	return nil
}
