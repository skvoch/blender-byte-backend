package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skvoch/blender-byte-backend/internal/application"
	localstore "github.com/skvoch/blender-byte-backend/internal/store/local"
)

func main() {
	logger := logrus.New()
	config := application.NewConfig()

	//if len(port) > 0 {
	//	config.BindPort = ":" + os.Getenv("PORT")
	//}

	store := localstore.NewLocalStore()

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
