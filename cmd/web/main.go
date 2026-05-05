package main

import (
	"net/http"

	"github.com/fujidaiti/bookings/internal/handlers"
	"github.com/fujidaiti/bookings/internal/middlewares"
	"github.com/fujidaiti/bookings/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const portNumber = ":8080"

func main() {
	err := repository.Init()
	if err != nil {
		panic(err)
	}
	defer repository.Db().Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger, middlewares.NoSurf)
	r.Get("/", handlers.Home)
	r.Get("/about", handlers.About)
	r.Get("/contact", handlers.Contact)
	r.Get("/generals-quarters", handlers.Generals)
	r.Get("/majors-suite", handlers.Majors)
	r.Get("/search", handlers.Search)
	r.Get("/booking", handlers.Booking)
	r.Post("/booking", handlers.PostBooking)
	r.Get("/reservation-summary", handlers.ReservationSummary)
	http.ListenAndServe(portNumber, r)
}
