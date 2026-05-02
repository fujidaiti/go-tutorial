package handlers

import (
	"net/http"

	"github.com/fujidaiti/bookings/pkg/renderer"
)

func Home(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "home")
}

func About(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "about")
}

func Contact(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "contact")
}

func Generals(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "generals")
}

func Majors(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "majors")
}

func SearchAvailability(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "search-availability")
}

func MakeReservation(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "make-reservation")
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary")
}
