package api

import (
	"encoding/json"
	ics "github.com/arran4/golang-ical"
	"github.com/purvisb179/profound-crow/internal/tasks"
	"github.com/purvisb179/profound-crow/pkg"
	"html/template"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type Handler struct {
	AsynqService *tasks.AsynqService
}

func NewHandler(asynqService *tasks.AsynqService) *Handler {
	return &Handler{AsynqService: asynqService}
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

	// Parse configuration JSON from the request
	configJSON := r.FormValue("configuration")
	var uploadInput pkg.UploadInput
	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &uploadInput); err != nil {
			http.Error(w, "Invalid configuration JSON", http.StatusBadRequest)
			log.Printf("Error parsing configuration JSON: %v", err)
			return
		}
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
		endTime, _ := event.GetEndAt()
		processTime = processTime.Local()
		summary := event.GetProperty(ics.ComponentPropertyDescription).Value

		payload := pkg.CalendarEventPayload{
			FilePath:      handler.Filename,
			EventSummary:  summary,
			EventStart:    processTime,
			EventEnd:      endTime,
			Configuration: uploadInput,
		}

		err := h.AsynqService.ProcessAndEnqueueCalendarEvent(payload)
		if err != nil {
			http.Error(w, "Error processing calendar event", http.StatusInternalServerError)
			log.Printf("error processing calendar event: %v", err)
			return
		}
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

	tasks, err := h.AsynqService.ListScheduledTasks()
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		log.Printf("Error retrieving tasks: %v", err)
		return
	}

	taskDetails := make([]pkg.CalendarTaskCheckResponse, len(tasks))

	for i, task := range tasks {
		var payload pkg.CalendarTaskPayload
		err := json.Unmarshal(task.Payload, &payload)
		if err != nil {
			http.Error(w, "Error unmarshalling payload", http.StatusInternalServerError)
			log.Printf("Error unmarshalling payload: %v", err)
			return
		}

		taskDetails[i] = pkg.CalendarTaskCheckResponse{
			Payload:   payload,
			StartTime: task.NextProcessAt,
		}
	}

	// Load the template from the given path
	funcMap := template.FuncMap{
		"formatTime": formatTime,
	}

	tmpl, err := template.New("check.html").Funcs(funcMap).ParseFiles("./web/check.html")

	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Error loading template: %v", err)
		return
	}

	// Render the template with the data
	err = tmpl.Execute(w, taskDetails)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Printf("Error rendering template: %v", err)
		return
	}
}

func formatTime(t time.Time) string {
	return t.Format("01/02/2006 03:04:05 PM")
}

func (h *Handler) ClearQueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract queue name from the request's query parameters.
	queueName := r.URL.Query().Get("queue")
	if queueName == "" {
		http.Error(w, "Queue name is required", http.StatusBadRequest)
		return
	}

	err := h.AsynqService.ClearQueue(queueName)
	if err != nil {
		http.Error(w, "Error clearing the queue", http.StatusInternalServerError)
		log.Printf("Error clearing the queue: %v", err)
		return
	}

	_, err = w.Write([]byte("Queue cleared successfully"))
	if err != nil {
		log.Printf("Error writing to buffer: %v", err)
	}
	return
}
