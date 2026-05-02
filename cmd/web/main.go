package main

import (
	"net/http"

	"github.com/fujidaiti/bookings/pkg/handlers"
)

const portNumber = ":8080"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)
	http.ListenAndServe(portNumber, nil)
}
