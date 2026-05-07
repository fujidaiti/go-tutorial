package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fujidaiti/bookings/internal/models"
	"github.com/fujidaiti/bookings/internal/renderer"
	"github.com/fujidaiti/bookings/internal/repository"
	"github.com/fujidaiti/bookings/internal/session"
	"golang.org/x/crypto/bcrypt"
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

func Standard(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "standard", renderer.DefaultData(r))
}

func Superior(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "superior", renderer.DefaultData(r))
}

func Deluxe(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "deluxe", renderer.DefaultData(r))
}

func Search(w http.ResponseWriter, r *http.Request) {
	handleSearch(w, r, "")
}

func SearchStandardRooms(w http.ResponseWriter, r *http.Request) {
	handleSearch(w, r, "Standard")
}

func SearchSuperiorRooms(w http.ResponseWriter, r *http.Request) {
	handleSearch(w, r, "Superior")
}

func SearchDeluxeRooms(w http.ResponseWriter, r *http.Request) {
	handleSearch(w, r, "Deluxe")
}

// TODO: Support pagination
func handleSearch(w http.ResponseWriter, r *http.Request, grade string) {
	q := r.URL.Query()
	data := renderer.DefaultData(r)
	data["Grade"] = strings.ToLower(grade)
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

	var gradeClause string
	args := []any{form.Start, form.End}
	if len(grade) > 0 {
		gradeClause = "AND g.name = $3"
		args = append(args, grade)
	}
	query := fmt.Sprintf(`
		SELECT r.id, r.name, g.name
		FROM rooms r
		JOIN grades g ON r.grade_id = g.id
		WHERE NOT EXISTS (
			SELECT 1
			FROM room_restrictions rr
			WHERE rr.room_id = r.id
				AND rr.arrival_date <= $2
				AND rr.departure_date >= $1
		)
		%s
		ORDER BY g.rank DESC;
	`, gradeClause)
	rows, err := repository.Db().Query(query, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type SearchResult struct {
		Room  models.Room
		Grade string
	}
	var results []SearchResult
	for rows.Next() {
		var r models.Room
		var g string
		if err := rows.Scan(&r.ID, &r.Name, &g); err != nil {
			panic(err)
		}
		results = append(results, SearchResult{r, g})
	}

	data["Results"] = results
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
	var grade string
	err := repository.Db().QueryRow(`
		SELECT r.id, r.name, g.name
		FROM rooms r
		JOIN grades g
		ON r.grade_id = g.id
		WHERE r.id = $1;
		`, id,
	).Scan(&room.ID, &room.Name, &grade)
	if err != nil {
		panic(err)
	}
	data["RoomId"] = room.ID
	data["RoomName"] = room.Name
	data["Grade"] = grade

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

// TODO: Make this page visible from only person who made this reservation.
func BookingDetails(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		panic(err)
	}

	var bk models.BookingForm
	var status string
	var roomName string
	var roomGrade string
	err = repository.Db().QueryRow(
		`SELECT
			r.name, b.arrival_date, b.departure_date, b.first_name,
			b.last_name, b.email, b.phone, b.status, g.Name
		FROM bookings b
		JOIN rooms r ON b.room_id = r.id
		JOIN grades g ON r.grade_id = g.id
		WHERE b.id = $1;
		`, id,
	).Scan(
		&roomName, &bk.Arrival, &bk.Departure, &bk.FirstName,
		&bk.LastName, &bk.Email, &bk.Phone, &status, &roomGrade,
	)
	if err != nil {
		panic(err)
	}

	if t, err := time.Parse(time.RFC3339, bk.Arrival); err == nil {
		bk.Arrival = t.Format("2006-01-02")
	}
	if t, err := time.Parse(time.RFC3339, bk.Departure); err == nil {
		bk.Departure = t.Format("2006-01-02")
	}

	data := renderer.DefaultData(r)
	data["RoomName"] = roomName
	data["RoomGrade"] = roomGrade
	data["Form"] = bk
	data["Status"] = status
	renderer.RenderTemplate(w, "booking-details", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
	renderer.RenderTemplate(w, "login-form", renderer.DefaultData(r))
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	data := renderer.DefaultData(r)

	form := models.LoginForm{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	result := form.Validate()
	data["Form"] = form
	data["FormResult"] = result
	if !result.IsValid() {
		renderer.RenderTemplate(w, "login-form", data)
		return
	}

	var id int
	var pwdHash string
	err := repository.Db().QueryRow(`
	SELECT g.id, g.pwd_hash
	FROM guests g
	WHERE g.email = $1
	`, form.Email,
	).Scan(&id, &pwdHash)

	if err == nil {
		// The email is already registered.
		err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(form.Password))
		if err == bcrypt.ErrMismatchedHashAndPassword {
			data["LoginErr"] = "The email and/or passward are incorrect."
			renderer.RenderTemplate(w, "login-form", data)
			return
		}
		if err != nil {
			panic(err)
		}
		// Successfully logged in.
		session.SetGuestCredential(id, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// The email is not yet registered, let's create a new account.
	formPwdHash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = repository.Db().QueryRow(`
		INSERT INTO guests (email, pwd_hash)
		VALUES ($1, $2)
		RETURNING id
		`, form.Email, formPwdHash,
	).Scan(&id)
	if err != nil {
		panic(err)
	}

	session.SetGuestCredential(id, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
