package store

import (
	model "github.com/skvoch/blender-byte-backend/internal/model"
)

// Store ...
type Store interface {
	RegisterUser(register *model.UserData) error
	Users() ([]*model.UserData, error)
	UserByLogin(login string) (*model.UserData, error)
}
