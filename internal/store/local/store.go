package localstore

import (
	model "github.com/skvoch/google-cloud-example/internal/model"
)

// LocalStore - In memory store
type LocalStore struct {
	users map[string]*model.UserData
}

// NewLocalStore ...
func NewLocalStore() *LocalStore {
	return &LocalStore{
		users: make(map[string]*model.UserData),
	}
}

// RegisterUser ...
func (l *LocalStore) RegisterUser(data *model.UserData) error {
	for _, user := range l.users {
		if user == data {
			return &model.UserAlreadyExistError{}
		}
	}

	l.users[data.Login] = data
	return nil
}

// Users ...
func (l *LocalStore) Users() ([]*model.UserData, error) {
	users := make([]*model.UserData, 0)

	for _, user := range l.users {
		users = append(users, user)
	}
	return users, nil
}

// UserByLogin ...
func (l *LocalStore) UserByLogin(login string) (*model.UserData, error) {
	if found := l.users[login]; found != nil {
		return found, nil
	}

	return nil, &model.CannotFindUserError{}
}
