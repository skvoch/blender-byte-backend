package store

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// FireStore ...
type FireStore struct {
	client *firestore.Client
	ctx    context.Context
}

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
