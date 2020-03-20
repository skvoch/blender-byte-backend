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

	a.router.HandleFunc("/v1.0/message/", a.handleMessage()).Methods("GET")
}

func (a *Application) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (a *Application) handleMessage() http.HandlerFunc {

	type Response struct {
		Text string `json:"id"`
	}

	type PythonResponse struct {
		RandomValue string `json:"RandomValue"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("http://127.0.0.1" + a.config.PythonPort + "/random/")

		if err != nil {
			a.logger.Error(err)
			return
		}
		pythonResonce := &PythonResponse{}
		if err := json.NewDecoder(res.Body).Decode(pythonResonce); err != nil {
			a.logger.Error(err)
			return
		}
		hash := a.getMD5Hash(pythonResonce.RandomValue)
		if err := a.store.AddHash(hash); err != nil {
			a.logger.Error(err)
			return
		}
		a.respond(w, r, http.StatusOK, &Response{
			Text: hash,
		})
	}
}

func (a *Application) getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
