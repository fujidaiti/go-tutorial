package models

import "regexp"

type LoginForm struct {
	Email    string
	Password string
}

type LoginFormValidationResult struct {
	EmailErr    string
	PasswordErr string
}

var rePwd = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func (f *LoginForm) Validate() LoginFormValidationResult {
	r := LoginFormValidationResult{}

	if len(f.Email) == 0 {
		r.EmailErr = "Email is required."
	}

	if len(f.Password) == 0 {
		r.PasswordErr = "Password is required."
	} else if !rePwd.MatchString(f.Password) {
		r.PasswordErr = "Password can contain only letters and numbers."
	}

	return r
}

func (r *LoginFormValidationResult) IsValid() bool {
	return len(r.EmailErr) == 0 && len(r.PasswordErr) == 0
}
