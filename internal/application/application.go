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

type ctxKey int8

const (
	ctxKeyRequestID ctxKey = iota
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
		config:  config,
		router:  mux.NewRouter(),
		logger:  logger,
		store:   store,
		logined: make(map[string]string),
	}
	application.setupHandlers()

	return application, nil
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *Application) setupHandlers() {
	a.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	a.router.HandleFunc("/v1.0/login/", a.handleLogin()).Methods("POST")
	a.router.HandleFunc("/v1.0/logout/", a.handleLogout()).Methods("POST")
	a.router.HandleFunc("/v1.0/register/", a.handleRegister()).Methods("POST")

	private := a.router.PathPrefix("/v1.0/private").Subrouter()
	private.Use(a.middlewareLogin)
	private.HandleFunc("/whoami/", a.handlePrivateWhoami()).Methods("GET")
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

func (a *Application) middlewareLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		if _, ok := a.logined[uuid]; ok {
			next.ServeHTTP(w, r)
		} else {
			a.respond(w, r, http.StatusNotFound, nil)
		}
	})
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
			a.error(w, r, http.StatusBadRequest, err)
		}

		a.respond(w, r, http.StatusCreated, nil)
	}
}

func (a *Application) handleLogin() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		request := &model.LoginRequest{}

		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			a.logger.Error(err)
			a.error(w, r, http.StatusBadRequest, err)
			return
		}

		user, err := a.store.UserByLogin(request.Login)

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		if user.Password == request.Password {
			uuid := uuid.New().String()

			a.logined[uuid] = user.Login
			a.respond(w, r, http.StatusOK, &model.LoginResponse{
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

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		if _, ok := a.logined[uuid]; ok {
			a.respond(w, r, http.StatusOK, &model.WhoamiResponse{
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
