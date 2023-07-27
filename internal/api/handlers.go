package api

import (
	ics "github.com/arran4/golang-ical"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"log"
	"net/http"
	"os"
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

	filePath := r.URL.Query().Get("filepath")
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Could not open file", http.StatusInternalServerError)
		log.Printf("could not open file: %v", err)
		return
	}
	defer file.Close()

	cal, err := ics.ParseCalendar(file)
	if err != nil {
		http.Error(w, "Could not parse calendar", http.StatusInternalServerError)
		log.Printf("could not parse calendar: %v", err)
		return
	}

	for _, event := range cal.Events() {
		processTime, _ := event.GetStartAt()
		summary := event.GetProperty(ics.ComponentPropertyDescription).Value

		task, err := tasks.CreateCalendarEvent(filePath, summary, processTime)
		if err != nil {
			http.Error(w, "Could not create task", http.StatusInternalServerError)
			log.Printf("could not create task: %v", err)
			return
		}

		durationUntilProcessing := processTime.Sub(time.Now())
		if durationUntilProcessing < 0 {
			log.Printf("event in the past, skipping: %v", err)
			continue
		}

		if _, err := h.Client.Enqueue(task, asynq.ProcessIn(durationUntilProcessing)); err != nil {
			http.Error(w, "Could not enqueue task", http.StatusInternalServerError)
			log.Printf("could not enqueue task: %v", err)
			return
		}
	}

	w.Write([]byte("Tasks created successfully"))
	return
}
