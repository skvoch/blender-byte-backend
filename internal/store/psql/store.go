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

// NewTest ...
func NewTest() (*PSQLStore, error) {
	db, err := gorm.Open("postgres", "host=34.77.221.9 port=5432 user=postgres dbname=test password=blender-byte")
	if err != nil {
		return nil, err
	}

	instance := &PSQLStore{
		db: db,
	}
	instance.applyMigrate()

	return instance, nil
}

func (p *PSQLStore) applyMigrate() {
	p.db.AutoMigrate(&model.UserData{})
}

// RegisterUser ...
func (p *PSQLStore) RegisterUser(data *model.UserData) error {
	errors := p.db.Create(data).GetErrors()

	for _, err := range errors {
		return err
	}

	return nil
}

// Users ...
func (p *PSQLStore) Users() ([]*model.UserData, error) {
	return nil, nil
}

// UserByLogin ...
func (p *PSQLStore) UserByLogin(login string) (*model.UserData, error) {
	user := &model.UserData{}

	errors := p.db.First(&user, "login = ?", login).GetErrors()

	for _, err := range errors {
		return nil, err
	}
	return user, nil
}
