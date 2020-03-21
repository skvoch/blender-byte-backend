package model

import "github.com/jinzhu/gorm"

// UserData ...
type UserData struct {
	gorm.Model
	Login    string `json:"login" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
}

// Book ...
type Book struct {
	Author      string  `json:"Author"`
	Code        string  `json:"Code"`
	Cost        float64 `json:"Cost"`
	Date        string  `json:"Date"`
	Description string  `json:"Description"`
	FullName    string  `json:"FullName"`
	ISBN        string  `json:"ISBN"`
	Name        string  `json:"Name"`
	Photo       string  `json:"Photo"`
	Publish     string  `json:"Publish"`
	Series      string  `json:"Series"`
	Sheets      int64   `json:"Sheets"`
	Topic       string  `json:"Topic"`
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
