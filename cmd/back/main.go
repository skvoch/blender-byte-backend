package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skvoch/blender-byte-backend/internal/application"
	psqlstore "github.com/skvoch/blender-byte-backend/internal/store/psql"
)

func main() {
	logger := logrus.New()
	config := application.NewConfig()

	store, err := psqlstore.New()

	if err != nil {
		logger.Error(err)
		return
	}

	app, err := application.New(store, config, logger)

	if err != nil {
		logger.Error(err)
		return
	}

	err = http.ListenAndServe(config.BindPort, app)

	if err != nil {
		logger.Error(err)
		return
	}
}
