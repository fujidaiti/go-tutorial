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
	http.ListenAndServe(portNumber, r)
}
