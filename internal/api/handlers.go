package api

import (
	ics "github.com/arran4/golang-ical"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"log"
	"mime/multipart"
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

	// max 10 MB file
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing the file", http.StatusInternalServerError)
		log.Printf("Error parsing the file: %v", err)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		log.Printf("Error retrieving the file: %v", err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error closing the file: %v", err)
		}
	}(file)

	log.Printf("Uploaded File: %+v", handler.Filename)
	log.Printf("File Size: %+v", handler.Size)
	log.Printf("MIME Header: %+v", handler.Header)

	cal, err := ics.ParseCalendar(file)
	if err != nil {
		http.Error(w, "Could not parse calendar", http.StatusInternalServerError)
		log.Printf("could not parse calendar: %v", err)
		return
	}

	for _, event := range cal.Events() {
		processTime, _ := event.GetStartAt()
		summary := event.GetProperty(ics.ComponentPropertyDescription).Value

		task, err := tasks.CreateCalendarEvent(handler.Filename, summary, processTime)
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

	_, err = w.Write([]byte("Tasks created successfully"))
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
	return

}
