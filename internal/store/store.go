package store

import (
	model "github.com/skvoch/blender-byte-backend/internal/model"
)

// Store ...
type Store interface {
	RegisterUser(register *model.UserData) error
	Users() ([]*model.UserData, error)
	UserByLogin(login string) (*model.UserData, error)

	AddType(_type *model.Type) (*model.Type, error)
	Types() ([]*model.Type, error)
	AssingBookToType(book *model.Book, _type *model.Type) error

	AddBook(book *model.Book) (*model.Book, error)
	BookIDsByType(typeID int) ([]uint, error)
	Book(ID uint) (model.Book, error)
}
