package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skvoch/google-cloud-example/internal/application"
)

func main() {
	logger := logrus.New()
	config := application.NewConfig()

	//if len(port) > 0 {
	//	config.BindPort = ":" + os.Getenv("PORT")
	//}

	app, err := application.New(config, logger, "blender-byte", "private.json")

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
