package handlers

import (
	"net/http"

	"github.com/fujidaiti/bookings/internal/models"
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

func Search(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "search", renderer.DefaultData(r))
}

func PostSearch(w http.ResponseWriter, r *http.Request) {
	form := models.SearchForm{
		Start: r.Form.Get("start"),
		End:   r.Form.Get("end"),
	}
	result := form.Validate()

	data := renderer.DefaultData(r)
	data["Form"] = form
	data["FormResult"] = result
	// TODO: Look up DB and return actual results
	if result.Valid() {
		data["IsFormValid"] = true
		data["AvailableRooms"] = []models.Room{
			{ID: 1, Name: "General's Quarters"},
			{ID: 2, Name: "Major's Suite"},
		}
	} else {
		data["IsFormValid"] = false
	}

	renderer.RenderTemplate(w, "search", data)
}

func MakeReservation(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "make-reservation", renderer.DefaultData(r))
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary", renderer.DefaultData(r))
}
