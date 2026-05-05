package handlers

import (
	"net/http"

	"github.com/fujidaiti/bookings/internal/models"
	"github.com/fujidaiti/bookings/internal/renderer"
	"github.com/fujidaiti/bookings/internal/repository"
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
	if !result.Valid() {
		data["IsFormValid"] = false
		renderer.RenderTemplate(w, "search", data)
		return
	} else {
		data["IsFormValid"] = true
	}

	rows, err := repository.Db().Query(`
		SELECT r.id, r.name
		FROM rooms r
		WHERE NOT EXISTS (
			SELECT 1
			FROM room_restrictions rr
			WHERE rr.room_id = r.id
				AND rr.arrival_date <= $2
				AND rr.departure_date >= $1
		);
	`, form.Start, form.End)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var r models.Room
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			panic(err)
		}
		rooms = append(rooms, r)
	}

	data["AvailableRooms"] = rooms
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

	renderer.RenderTemplate(w, "booking-form", data)
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
	if !result.Valid() {
		data["RoomId"] = id
		data["RoomName"] = "Major's quarter"
		data["Form"] = form
		data["FormResult"] = result
		renderer.RenderTemplate(w, "booking-form", data)
		return
	}

	// TODO: Check if the room is available
	// TODO: Save the reservation data to DB
	data["Form"] = form
	data["RoomName"] = "Major's quarter"
	renderer.RenderTemplate(w, "booking-summary", data)
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary", renderer.DefaultData(r))
}
