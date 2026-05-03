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
	q := r.URL.Query()
	data := renderer.DefaultData(r)
	if !q.Has("start") && !q.Has("end") {
		renderer.RenderTemplate(w, "search", data)
		return
	}

	form := models.SearchForm{
		Start: q.Get("start"),
		End:   q.Get("end"),
	}
	result := form.Validate()

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

func Book(w http.ResponseWriter, r *http.Request) {
	data := renderer.DefaultData(r)
	data["FormResult"] = models.BookingFormValidationResult{}

	id := r.PathValue("roomId")
	// TODO: Get the actual room name from DB
	room := "Major's qurter"
	data["RoomId"] = id
	data["RoomName"] = room

	q := r.URL.Query()
	los := models.SearchForm{
		Start: q.Get("arrival"),
		End:   q.Get("departure"),
	}
	result := los.Validate()
	if result.Valid() {
		data["Form"] = models.BookingForm{
			Arrival:   los.Start,
			Departure: los.End,
		}
	} else {
		data["Form"] = models.BookingForm{}
	}

	renderer.RenderTemplate(w, "book", data)
}

func PostBook(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if ID is valid
	id := r.PathValue("roomId")

	form := models.BookingForm{
		Arrival:   r.Form.Get("arrival"),
		Departure: r.Form.Get("departure"),
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	result := form.Validate()

	data := renderer.DefaultData(r)
	if result.Valid() {
		// TODO: Check if the room is available
		// TODO: Save the reservation data to DB
		// TODO: Redirect to the result page
		w.Write([]byte("The room has been booked for you!"))
	} else {
		data["RoomId"] = id
		data["RoomName"] = "Major's quarter"
		data["Form"] = form
		data["FormResult"] = result
		renderer.RenderTemplate(w, "book", data)
	}
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary", renderer.DefaultData(r))
}
