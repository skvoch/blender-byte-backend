package psqlstore

import (
	"github.com/jinzhu/gorm"
	"github.com/skvoch/blender-byte-backend/internal/model"

	// ...
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PSQLStore ...
type PSQLStore struct {
	db *gorm.DB
}

// New - heper function
func New() (*PSQLStore, error) {
	db, err := gorm.Open("postgres", "host=34.77.221.9 port=5432 user=postgres dbname=dev password=blender-byte")
	if err != nil {
		return nil, err
	}

	return &PSQLStore{
		db: db,
	}, nil
}

// RegisterUser ...
func (p *PSQLStore) RegisterUser(data *model.UserData) error {
	return nil
}

// Users ...
func (p *PSQLStore) Users() ([]*model.UserData, error) {
	return nil, nil
}

// UserByLogin ...
func (p *PSQLStore) UserByLogin(login string) (*model.UserData, error) {
	return nil, nil
}
