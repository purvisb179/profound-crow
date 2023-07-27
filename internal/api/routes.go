package api

import (
	"github.com/go-chi/chi"
)

func BindRoutes(r *chi.Mux) {
	r.Post("/new-calendar", CreateCalendarHandler)
}
