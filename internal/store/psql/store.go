package psqlstore

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/skvoch/blender-byte-backend/internal/model"
	gormbulk "github.com/t-tiger/gorm-bulk-insert"

	// ...
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PSQLStore ...
type PSQLStore struct {
	db *gorm.DB
}

// New - heper function
func New() (*PSQLStore, error) {

	dbString := ""

	if len(os.Getenv("ENV_CLOUD")) > 0 {
		dbString = "host=/cloudsql/blender-byte:europe-west1:sql-db port=5432 user=postgres dbname=dev password=blender-byte"
	} else {
		dbString = "host=34.77.221.9 port=5432 user=postgres dbname=dev password=blender-byte"
	}

	db, err := gorm.Open("postgres", dbString)
	if err != nil {
		return nil, err
	}

	instance := &PSQLStore{
		db: db,
	}

	instance.applyMigrate()
	return instance, nil
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
	p.db.AutoMigrate(&model.Book{})
	p.db.AutoMigrate(&model.Type{})
}

// Clean ...
func (p *PSQLStore) Clean() {
	p.db.DB().Query("TRUNCATE TABLE user_data;")
	p.db.DB().Query("TRUNCATE TABLE types;")
	p.db.DB().Query("TRUNCATE TABLE books;")
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

// AddType ...
func (p *PSQLStore) AddType(_type *model.Type) (*model.Type, error) {
	errors := p.db.Create(_type).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return _type, nil
}

// Types ...
func (p *PSQLStore) Types() ([]*model.Type, error) {
	types := make([]*model.Type, 0)

	errors := p.db.Find(&types).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return types, nil
}

// AssingBookToType ...
func (p *PSQLStore) AssingBookToType(book *model.Book, _type *model.Type) error {
	_type.Books = append(_type.Books, *book)
	errors := p.db.Save(_type).GetErrors()

	for _, err := range errors {
		return err
	}

	return nil
}

// AddBook ...
func (p *PSQLStore) AddBook(book *model.Book) (*model.Book, error) {
	errors := p.db.Create(book).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return book, nil
}

// AddBooks ...
func (p *PSQLStore) AddBooks(books []*model.Book) ([]*model.Book, error) {
	var insertRecords []interface{}

	for _, book := range books {
		insertRecords = append(insertRecords, book)
	}

	gormbulk.BulkInsert(p.db, insertRecords, 3000)

	return books, nil
}

// BookIDsByType ...
func (p *PSQLStore) BookIDsByType(typeID int) ([]uint, error) {
	books := make([]*model.Book, 0)

	errors := p.db.Find(&books, "type_id = ?", typeID).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	result := make([]uint, 0)

	for _, book := range books {
		result = append(result, book.ID)
	}
	return result, nil
}

// BooksByType ...
func (p *PSQLStore) BooksByType(typeID int) ([]*model.Book, error) {
	books := make([]*model.Book, 0)

	errors := p.db.Find(&books, "type_id = ?", typeID).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return books, nil
}

// Book ...
func (p *PSQLStore) Book(ID uint) (*model.Book, error) {
	book := &model.Book{}

	errors := p.db.First(&book, "id = ?", ID).GetErrors()

	for _, err := range errors {
		return nil, err
	}
	return book, nil
}

// FindBook ...
func (p *PSQLStore) FindBook(key string) ([]*model.Book, error) {

	books := make([]*model.Book, 0)

	errors := p.db.Where("position(? in full_name) > 0", key).Find(&books).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return books, nil
}

// FindBookByTag ...
func (p *PSQLStore) FindBookByTag(tag string) ([]uint, error) {

	books := make([]*model.Book, 0)

	errors := p.db.Where("position(? in tags) > 0", tag).Find(&books).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	result := make([]uint, 0)

	for _, book := range books {
		result = append(result, book.ID)
	}

	return result, nil
}

// AddTag ...
func (p *PSQLStore) AddTag(tag *model.Tag) (*model.Tag, error) {
	errors := p.db.Create(tag).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return tag, nil
}

// Tags ...
func (p *PSQLStore) Tags(tag *model.Tag) ([]*model.Tag, error) {
	tags := make([]*model.Tag, 0)

	errors := p.db.Find(&tags).GetErrors()

	for _, err := range errors {
		return nil, err
	}

	return tags, nil
}
