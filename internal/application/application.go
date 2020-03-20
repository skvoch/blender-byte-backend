package application

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	store "github.com/skvoch/google-cloud-example/internal/firestore"
	model "github.com/skvoch/google-cloud-example/internal/model"
)

// Application - REST backend
type Application struct {
	config *Config
	router *mux.Router
	logger *logrus.Logger
	store  *store.FireStore
}

// New - helper function
func New(config *Config, logger *logrus.Logger, name string, jsonPath string) (*Application, error) {
	store, err := store.New(name, jsonPath)

	if err != nil {
		return nil, err
	}

	application := &Application{
		config: config,
		router: mux.NewRouter(),
		logger: logger,
		store:  store,
	}
	application.setupHandlers()

	return application, nil
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *Application) setupHandlers() {
	a.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	//a.router.HandleFunc("/v1.0/message/", a.handleMessage()).Methods("GET")
}

func (a *Application) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	a.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (a *Application) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (a *Application) handleRegister() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		registerData := &model.RegisterRequest{}

		if err := json.NewDecoder(r.Body).Decode(registerData); err != nil {
			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		if state := registerData.Validate(); state == false {
			err := &model.FailedValidationError{}

			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}
	}
}

func (a *Application) getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
