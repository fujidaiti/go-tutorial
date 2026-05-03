package models

import "time"

type BookingForm struct {
	Arrival   string
	Departure string
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

type BookingFormValidationResult struct {
	ArrivalErr   string
	DepartureErr string
	FirstNameErr string
	LastNameErr  string
	EmailErr     string
	PhoneErr     string
}

func (r *BookingFormValidationResult) Valid() bool {
	return len(r.ArrivalErr) == 0 &&
		len(r.DepartureErr) == 0 &&
		len(r.FirstNameErr) == 0 &&
		len(r.LastNameErr) == 0 &&
		len(r.EmailErr) == 0 &&
		len(r.PhoneErr) == 0
}

func (f *BookingForm) Validate() BookingFormValidationResult {
	result := BookingFormValidationResult{}

	if len(f.Arrival) == 0 {
		result.ArrivalErr = "Arrival date is required."
	}
	if len(f.Departure) == 0 {
		result.DepartureErr = "Departure date is required."
	}
	if len(f.FirstName) == 0 {
		result.FirstNameErr = "Fist name is required."
	}
	if len(f.LastName) == 0 {
		result.LastNameErr = "Last name is required."
	}
	if len(f.Email) == 0 {
		result.EmailErr = "Email is required."
	}

	arr, arrErr := time.Parse(dateFormat, f.Arrival)
	if arrErr != nil {
		result.ArrivalErr = "Invalid arrival date format. Please use YYYY-MM-DD."
	}
	dep, depErr := time.Parse(dateFormat, f.Departure)
	if depErr != nil {
		result.ArrivalErr = "Invalid departure date format. Please use YYYY-MM-DD."
	}
	if arrErr == nil && depErr == nil && dep.Before(arr) {
		result.DepartureErr = "Departure date must be after arrival date."
	}

	return result
}
