package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
)

// FireStore ...
type FireStore struct {
	client *firestore.Client
	ctx    context.Context
}

/*
// AddHash ...
func (f *FireStore) AddHash(hash string) error {
	_, _, err := f.client.Collection("hashes").Add(f.ctx, map[string]interface{}{
		"hash": hash,
	})

	if err != nil {
		return err
	}

	return nil
}

// RegisterUser ...
func (f *FireStore) RegisterUser(data *model.UserData) error {
	/*users := f.client.Collection("users")


	_, _, err := f.client.Collection("users").Add(f.ctx, map[string]interface{}{
		"hash": hash,
	})
}

// Users ...
func (f *FireStore) Users() ([]*model.UserData, error) {

}

// UserByLogin ...
func (f *FireStore) UserByLogin(login string) (*model.UserData, error) {

}

// New - helper function
func New(name string, jsonPath string) (*FireStore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, name, option.WithCredentialsFile(jsonPath))

	if err != nil {
		return nil, err
	}
	return &FireStore{
		client: client,
		ctx:    ctx,
	}, nil
}
*/
