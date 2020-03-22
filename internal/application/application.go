package application

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	model "github.com/skvoch/blender-byte-backend/internal/model"
	psqlstore "github.com/skvoch/blender-byte-backend/internal/store/psql"
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

	store *psqlstore.PSQLStore

	logined    map[string]string
	loginedMux sync.Mutex
}

//, name string, jsonPath string

// New - helper function
func New(store *psqlstore.PSQLStore, config *Config, logger *logrus.Logger) (*Application, error) {

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

	a.router.HandleFunc("/v1.0/types/", a.handleTypes()).Methods("GET")
	a.router.HandleFunc("/v1.0/types/{id}/", a.handleBooksIDs()).Methods("GET")
	a.router.HandleFunc("/v1.0/types/{id}/count/", a.handleTypeBooksCount()).Methods("GET")
	a.router.HandleFunc("/v1.0/types/{id}/books/", a.handleTypeBooks()).Methods("GET")

	a.router.HandleFunc("/v1.0/books/{id}/", a.handleBookdByID()).Methods("GET")
	a.router.HandleFunc("/v1.0/find/", a.handleFind()).Methods("GET")
	a.router.HandleFunc("/v1.0/find_tag/", a.handleFindTag()).Methods("GET")
	a.router.HandleFunc("/v1.0/tags/", a.handleTags()).Methods("GET")

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

func (a *Application) uuidIsLogined(uuid string) bool {
	result := false

	a.loginedMux.Lock()
	if _, ok := a.logined[uuid]; ok {
		result = ok
	}
	a.loginedMux.Unlock()

	return result
}

func (a *Application) uuidRemoveLogined(uuid string) {
	a.loginedMux.Lock()

	a.logined[uuid] = ""
	a.loginedMux.Unlock()
}

func (a *Application) uuidSetLogined(uuid string, login string) {
	a.loginedMux.Lock()
	a.logined[uuid] = login
	a.loginedMux.Unlock()
}

func (a *Application) uuidGetLogin(uuid string) string {
	result := ""

	if a.uuidIsLogined(uuid) {
		a.loginedMux.Lock()
		result = a.logined[uuid]
		a.loginedMux.Unlock()
	}

	return result
}

func (a *Application) middlewareLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		if a.uuidIsLogined(uuid) {
			next.ServeHTTP(w, r)
		} else {
			a.respond(w, r, http.StatusNotFound, nil)
		}
	})
}

func (a *Application) handleTags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := a.store.Tags()

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		a.respond(w, r, http.StatusOK, tags)
	}
}

func (a *Application) handleFindTag() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		books, err := a.store.FindBookByTag(key)

		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.respond(w, r, http.StatusOK, books)
	}
}

func (a *Application) handleFind() http.HandlerFunc {

	type Response struct {
		BooksIDs []uint `json:"books_ids"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		books, err := a.store.FindBook(key)

		if err != nil {
			a.error(w, r, http.StatusInternalServerError, err)
			return
		}

		a.respond(w, r, http.StatusOK, books)

	}
}

func (a *Application) handleTypeBooks() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, err := strconv.Atoi(vars["id"])

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		books, err := a.store.BooksByType(ID)

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		count, err1 := strconv.Atoi(r.URL.Query().Get("count"))
		page, err2 := strconv.Atoi(r.URL.Query().Get("page"))

		if err1 != nil || err2 != nil || (count*page > len(books)) {

			a.respond(w, r, http.StatusOK, books)
		} else {
			a.respond(w, r, http.StatusOK, books[count*page:(count*page)+count])
		}
	}

}

func (a *Application) handleBooksIDs() http.HandlerFunc {

	type Response struct {
		BooksIDs []uint `json:"books_ids"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, err := strconv.Atoi(vars["id"])

		IDs, err := a.store.BookIDsByType(ID)

		if err != nil {
			a.error(w, r, http.StatusNotFound, err)
			return
		}

		count, err1 := strconv.Atoi(r.URL.Query().Get("count"))
		page, err2 := strconv.Atoi(r.URL.Query().Get("page"))

		if err1 != nil || err2 != nil || (count*page > len(IDs)) {
			res := &Response{
				BooksIDs: IDs,
			}
			a.respond(w, r, http.StatusOK, res)
		} else {
			res := &Response{
				BooksIDs: IDs[count*page : (count*page)+count],
			}
			a.respond(w, r, http.StatusOK, res)
		}
	}
}

func (a *Application) handleBookdByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, err := strconv.Atoi(vars["id"])

		book, err := a.store.Book((uint)(ID))

		if err != nil {
			a.error(w, r, http.StatusNotFound, err)
			return
		}

		a.respond(w, r, http.StatusOK, book)
	}
}

func (a *Application) handleTypeBooksCount() http.HandlerFunc {

	type Response struct {
		Count int `json:"count"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ID, err := strconv.Atoi(vars["id"])

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		IDs, err := a.store.BookIDsByType(ID)

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		a.respond(w, r, http.StatusOK, &Response{
			Count: len(IDs),
		})
	}
}

func (a *Application) handleTypes() http.HandlerFunc {

	type Response struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		types, err := a.store.Types()

		if err != nil {
			a.error(w, r, http.StatusBadRequest, err)
		}

		res := make([]*Response, 0)

		for _, t := range types {
			res = append(res, &Response{
				Name: t.Name,
				ID:   t.ID,
			})
		}
		a.respond(w, r, http.StatusOK, res)
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

			a.uuidSetLogined(uuid, user.Login)

			a.respond(w, r, http.StatusOK, &model.LoginResponse{
				PrivateUUID: uuid,
			})

		}
	}
}

func (a *Application) handleLogout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		a.uuidRemoveLogined(uuid)

		a.respond(w, r, http.StatusOK, nil)
	}
}

func (a *Application) handlePrivateWhoami() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		uuid := r.Header.Get("X-PRIVATE-UUID")

		if a.uuidIsLogined(uuid) {
			a.respond(w, r, http.StatusOK, &model.WhoamiResponse{
				Login: a.uuidGetLogin(uuid),
			})
		}
	}
}

func (a *Application) getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
