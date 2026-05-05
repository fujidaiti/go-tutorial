package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

func Booking(w http.ResponseWriter, r *http.Request) {
	data := renderer.DefaultData(r)
	data["FormResult"] = models.BookingFormValidationResult{}

	q := r.URL.Query()
	id := q.Get("roomId")
	if len(id) == 0 {
		panic("roomId is required.")
	}

	var room models.Room
	err := repository.Db().QueryRow(
		"SELECT id, name FROM rooms r WHERE id = $1;", id,
	).Scan(&room.ID, &room.Name)
	if err != nil {
		panic(err)
	}
	data["RoomId"] = room.ID
	data["RoomName"] = room.Name

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

func PostBooking(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		panic(err)
	}

	form := models.BookingForm{
		Arrival:   r.Form.Get("arrival"),
		Departure: r.Form.Get("departure"),
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	result := form.Validate()

	var room models.Room
	err = repository.Db().QueryRow(
		"SELECT id, name FROM rooms WHERE id = $1;", roomId,
	).Scan(&room.ID, &room.Name)
	if err != nil {
		panic(err)
	}

	data := renderer.DefaultData(r)
	data["RoomId"] = room.ID
	data["RoomName"] = room.Name
	data["Form"] = form
	data["FormResult"] = result

	if result.Valid() {
		data["FormValid"] = true
	} else {
		data["FormValid"] = false
		renderer.RenderTemplate(w, "booking-form", data)
		return
	}

	var available bool
	err = repository.Db().QueryRow(
		`SELECT NOT EXISTS (
			SELECT 1
			FROM room_restrictions rr
			WHERE rr.room_id = $1
				AND rr.arrival_date <= $3
				AND rr.departure_date >= $2
			);
		`,
		roomId, form.Arrival, form.Departure,
	).Scan(&available)
	if err != nil {
		panic(err)
	}
	if available {
		data["RoomAvailable"] = true
	} else {
		data["RoomAvailable"] = false
		renderer.RenderTemplate(w, "booking-form", data)
		return
	}

	tx, err := repository.Db().Begin()
	if err != nil {
		panic(err)
	}

	var bookingId int
	err = tx.QueryRow(
		`INSERT INTO bookings (
			first_name, last_name, email, phone,
			arrival_date, departure_date, room_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;
		`,
		form.FirstName, form.LastName, form.Email, form.Phone,
		form.Arrival, form.Departure, roomId,
	).Scan(&bookingId)
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(
		`INSERT INTO room_restrictions (
			arrival_date, departure_date, room_id, booking_id
		)
		VALUES ($1, $2, $3, $4);
		`,
		form.Arrival, form.Departure, roomId, bookingId,
	)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, fmt.Sprintf("/booking/%d", bookingId), 303)
}

func BookingDetails(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "booking-details", renderer.DefaultData(r))
}

func ReservationSummary(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "reservation-summary", renderer.DefaultData(r))
}
