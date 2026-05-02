package main

import (
	"net/http"

	"github.com/fujidaiti/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const portNumber = ":8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", handlers.Home)
	r.Get("/about", handlers.About)
	r.Get("/contact", handlers.Contact)
	r.Get("/generals-quarters", handlers.Generals)
	r.Get("/majors-suite", handlers.Majors)
	r.Get("/search-availability", handlers.SearchAvailability)
	r.Get("/make-reservation", handlers.MakeReservation)
	r.Get("/reservation-summary", handlers.ReservationSummary)
	http.ListenAndServe(portNumber, r)
}
