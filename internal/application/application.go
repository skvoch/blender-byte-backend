package application

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	model "github.com/skvoch/blender-byte-backend/internal/model"
	store "github.com/skvoch/blender-byte-backend/internal/store"
)

// Application - REST backend
type Application struct {
	config *Config
	router *mux.Router
	logger *logrus.Logger

	store store.Store

	logined map[string]string
}

//, name string, jsonPath string

// New - helper function
func New(store store.Store, config *Config, logger *logrus.Logger) (*Application, error) {

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
		userData := &model.UserData{}

		if err := json.NewDecoder(r.Body).Decode(userData); err != nil {
			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		if state := userData.Validate(); state == false {
			err := &model.FailedValidationError{}

			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := a.store.RegisterUser(userData); err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (a *Application) handleLogin() http.HandlerFunc {

	type LoginRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	type Response struct {
		PrivateUUID string `json:"private_uuid"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := &LoginRequest{}

		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := a.store.UserByLogin(request.Login)

		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
		}

		if user.Password == request.Password {
			uuid := uuid.New().String()

			a.logined[uuid] = user.Login
			a.respond(w, r, http.StatusOK, &Response{
				PrivateUUID: uuid,
			})
		}
	}
}

func (a *Application) handleLogout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		a.logined[uuid] = ""

		a.respond(w, r, http.StatusOK, nil)
	}
}

func (a *Application) handlePrivateWhoami() http.HandlerFunc {

	type Response struct {
		Login string `json:"login"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		if a.logined[uuid] != "" {
			a.respond(w, r, http.StatusOK, &Response{
				Login: a.logined[uuid],
			})
		}
	}
}

func (a *Application) getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
