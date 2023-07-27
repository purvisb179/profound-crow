package api

import (
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"log"
	"net/http"
	"time"
)

var client = asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"}) //TODO dont make new client on every request

func CreateCalendarHandler(w http.ResponseWriter, r *http.Request) {
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

	if _, err := client.Enqueue(task); err != nil {
		http.Error(w, "Error enqueuing task", http.StatusInternalServerError)
		log.Printf("Error enqueuing task: %v", err)
		return
	}

	log.Printf("Enqueued task: %v", task)
	w.Write([]byte("Task created successfully"))
	return
}
