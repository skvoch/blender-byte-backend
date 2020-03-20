package store 

import (
	model "github.com/skvoch/google-cloud-example/internal/model"
)

// Store ...
type Store interface {
	RegisterUser(register *model.UserData) error
	Users() ([]*model.UserData, error)
	UserByLogin(login string) ( *model.UserData,error)
}