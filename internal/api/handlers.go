package api

import (
	"log"
	"net/http"
)

func CreateCalendarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("received")
	return
}
