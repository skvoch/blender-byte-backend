package application

import (
	"github.com/sirupsen/logrus"
	psqlstore "github.com/skvoch/blender-byte-backend/internal/store/psql"
)

// NewTestApplication - helper func
func NewTestApplication() (*Application, error) {
	store, err := psqlstore.NewTest()
	store.Clean()

	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	config := NewConfig()

	app, err := New(store, config, logger)

	if err != nil {
		return nil, err
	}

	return app, nil
}
