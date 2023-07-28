package api

import (
	"encoding/json"
	ics "github.com/arran4/golang-ical"
	"github.com/hibiken/asynq"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type Handler struct {
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewHandler(client *asynq.Client, inspector *asynq.Inspector) *Handler {
	return &Handler{Client: client, Inspector: inspector}
}

func (h *Handler) CreateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

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

	log.Printf("Calendar parsed successfully, processing events.")

	for _, event := range cal.Events() {
		processTime, _ := event.GetStartAt()
		processTime = processTime.Local() // convert to local timezone

		summary := event.GetProperty(ics.ComponentPropertyDescription).Value

		task, err := tasks.CreateCalendarEvent(handler.Filename, summary, processTime)
		if err != nil {
			http.Error(w, "Could not create task", http.StatusInternalServerError)
			log.Printf("could not create task: %v", err)
			return
		}

		log.Printf("Task created successfully. Summary: %s, Process Time: %v", summary, processTime)

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

		log.Printf("Task enqueued successfully. Duration until processing: %v", durationUntilProcessing)
	}

	_, err = w.Write([]byte("Tasks created successfully"))
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
	return
}

func (h *Handler) CheckQueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.Inspector.ListScheduledTasks("default", 0, -1)
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		log.Printf("Error retrieving tasks: %v", err)
		return
	}

	taskDetails := make([]map[string]interface{}, len(tasks))

	for i, task := range tasks {
		taskDetails[i] = map[string]interface{}{
			"ID":      task.ID,
			"Type":    task.Type,
			"Payload": task.Payload,
		}
	}

	response, err := json.Marshal(taskDetails)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		log.Printf("Error marshalling response: %v", err)
		return
	}

	w.Write(response)
	return
}
