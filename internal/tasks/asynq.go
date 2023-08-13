package tasks

import (
	"github.com/hibiken/asynq"
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
