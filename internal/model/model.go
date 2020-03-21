package model

import "github.com/jinzhu/gorm"

// UserData ...
type UserData struct {
	gorm.Model

	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// NewTestUser - helper func
func NewTestUser() *UserData {
	return &UserData{
		Login:    "ExampleLogin",
		Password: "ExamplePassword",
		Email:    "example@email.com",
	}
}

// Validate ...
func (r *UserData) Validate() bool {
	if len(r.Email) < 6 {
		return false
	}

	if len(r.Login) < 6 {
		return false
	}

	if len(r.Password) < 6 {
		return false
	}

	return true
}
