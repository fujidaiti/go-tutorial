package handlers

import (
	"fmt"
	"net/http"

	"github.com/fujidaiti/bookings/internal/renderer"
)

func Home(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "home", renderer.DefaultData(r))
}

func About(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "about", renderer.DefaultData(r))
}

func Contact(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "contact", renderer.DefaultData(r))
}

func Generals(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "generals", renderer.DefaultData(r))
}

func Majors(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "majors", renderer.DefaultData(r))
}

func SearchAvailability(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "search-availability", renderer.DefaultData(r))
}

func PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	fmt.Printf("Start: %s, End: %s\n", start, end)

	data := renderer.DefaultData(r)
	data["Start"] = start
	data["End"] = end
	renderer.RenderTemplate(w, "no-room-available", data)
}

func MakeReservation(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "make-reservation", renderer.DefaultData(r))
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary", renderer.DefaultData(r))
}
