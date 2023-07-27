package api

import (
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	Client *asynq.Client
}

func NewHandler(client *asynq.Client) *Handler {
	return &Handler{Client: client}
}

func (h *Handler) CreateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	task, err := tasks.CreateCalendarEvent("dummy/path", "dummy event", time.Now())
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		log.Printf("Error creating task: %v", err)
		return
	}

	if _, err := h.Client.Enqueue(task); err != nil {
		http.Error(w, "Error enqueuing task", http.StatusInternalServerError)
		log.Printf("Error enqueuing task: %v", err)
		return
	}

	log.Printf("Enqueued task: %v", task)
	w.Write([]byte("Task created successfully"))
	return
}
