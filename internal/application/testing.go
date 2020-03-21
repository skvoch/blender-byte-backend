package application

import (
	"github.com/sirupsen/logrus"
	localstore "github.com/skvoch/blender-byte-backend/internal/store/local"
)

// NewTestApplication - helper func
func NewTestApplication() (*Application, error) {
	store := localstore.New()
	logger := logrus.New()
	config := NewConfig()

	app, err := New(store, config, logger)

	if err != nil {
		return nil, err
	}

	return app, nil
}
