package models

type LoginFormData struct {
	Email    string
	Password string
	Errors   map[string]error
}

type SignupFormData struct {
	Name     string
	Email    string
	Password string
	Errors   map[string]error
}
