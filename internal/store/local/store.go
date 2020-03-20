package localstore


import (
	model "github.com/skvoch/google-cloud-example/internal/model"
)


// LocalStore - In memory store
type LocalStore struct {
	users map[string]*model.UserData
}

// NewLocalStore ...
func NewLocalStore() *LocalStore{
	return &LocalStore{
		users: make(map[string]*model.UserData),
	}
}

// RegisterUser ...
func (l *LocalStore)RegisterUser(data *model.UserData) error {
	for _, user := range l.users {
		if user == data {
			return err
		}
	}

	l.users[data.Login] := user
	return nil
}

// Users ...
func (l* LocalStore)Users() ([]*model.UserData, error) {

}

// UserByLogin ...
func (l *LocalStore)UserByLogin(login string) (*model.UserData,error) {

}
