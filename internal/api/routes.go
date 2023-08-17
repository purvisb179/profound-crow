package api

import (
	"github.com/go-chi/chi"
	"net/http"
)

func BindRoutes(r *chi.Mux, apiHandler *Handler) {
	r.Post("/new-calendar", func(w http.ResponseWriter, r *http.Request) {
		apiHandler.CreateCalendarHandler(w, r)
	})
	r.Get("/check-queue", apiHandler.CheckQueueHandler)
	r.Post("/clear-queue", apiHandler.ClearQueueHandler)

}
