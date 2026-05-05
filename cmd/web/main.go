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
	r.Get("/standard", handlers.Standard)
	r.Get("/standard/search", handlers.SearchStandardRooms)
	r.Get("/superior", handlers.Superior)
	r.Get("/superior/search", handlers.SearchSuperiorRooms)
	r.Get("/deluxe", handlers.Deluxe)
	r.Get("/deluxe/search", handlers.SearchDeluxeRooms)
	r.Get("/search", handlers.Search)
	r.Get("/booking", handlers.Booking)
	r.Post("/booking", handlers.PostBooking)
	r.Get("/booking/{id}", handlers.BookingDetails)
	http.ListenAndServe(portNumber, r)
}
