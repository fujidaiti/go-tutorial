package models

import "time"

type SearchForm struct {
	Start string
	End   string
}

type SearchFormValidationResult struct {
	HasStartDateErr bool
	HasEndDateErr   bool
	StartDateErr    string
	EndDateErr      string
}

const dateFormat = "2006-01-02"

func (s *SearchFormValidationResult) Valid() bool {
	return !s.HasStartDateErr && !s.HasEndDateErr
}

func (s *SearchForm) Validate() SearchFormValidationResult {
	result := SearchFormValidationResult{}

	if len(s.Start) == 0 {
		result.HasStartDateErr = true
		result.StartDateErr = "Start date is required."
	}
	if len(s.End) == 0 {
		result.HasEndDateErr = true
		result.EndDateErr = "End date is required."
	}
	if !result.Valid() {
		return result
	}

	start, err := time.Parse(dateFormat, s.Start)
	if err != nil {
		result.HasStartDateErr = true
		result.StartDateErr = "Invalid start date format. Please use YYYY-MM-DD."
		return result
	}

	end, err := time.Parse(dateFormat, s.End)
	if err != nil {
		result.HasEndDateErr = true
		result.EndDateErr = "Invalid end date format. Please use YYYY-MM-DD."
		return result
	}

	if end.Before(start) {
		result.HasEndDateErr = true
		result.EndDateErr = "End date must be after start date."
		return result
	}

	return result
}
