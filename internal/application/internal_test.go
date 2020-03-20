package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skvoch/blender-byte-backend/internal/application"
	"github.com/skvoch/blender-byte-backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	app, err := application.NewTestApplication()

	assert.NoError(t, err)
	assert.NotNil(t, app)

	cases := []struct {
		Name     string
		UserData *model.UserData
		Error    bool
	}{
		{
			Name:     "Empty",
			UserData: &model.UserData{},
			Error:    true,
		},
		{
			Name:     "Valid",
			UserData: model.NewTestUser(),
			Error:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			json, err := json.Marshal(c.UserData)
			assert.NoError(t, err)

			reader := bytes.NewReader(json)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/v1.0/register/", reader)

			app.ServeHTTP(rec, req)

			if c.Error == true {
				assert.Equal(t, http.StatusBadRequest, rec.Code)

			} else {
				assert.Equal(t, http.StatusCreated, rec.Code)
			}
		})
	}

}

func TestPrivate(t *testing.T) {
	app, err := application.NewTestApplication()

	assert.NoError(t, err)
	assert.NotNil(t, app)

	user := model.NewTestUser()

	js, err := json.Marshal(user)
	assert.NoError(t, err)

	reader := bytes.NewReader(js)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1.0/register/", reader)

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Try to get private without login

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1.0/private/whoami/", reader)

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Try login
	js, err = json.Marshal(&model.LoginRequest{
		Login:    user.Login,
		Password: user.Password,
	})
	assert.NoError(t, err)

	reader = bytes.NewReader(js)

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/v1.0/login/", reader)

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	loginResponse := &model.LoginResponse{}
	json.Unmarshal(rec.Body.Bytes(), loginResponse)

	// Try get private with private token

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1.0/private/whoami/", reader)
	req.Header.Set("X-PRIVATE-UUID", loginResponse.PrivateUUID)

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	whoamiResponse := &model.WhoamiResponse{}
	if err := json.Unmarshal(rec.Body.Bytes(), whoamiResponse); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, whoamiResponse.Login, user.Login)

	// Logout

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/v1.0/logout/", reader)
	req.Header.Set("X-PRIVATE-UUID", loginResponse.PrivateUUID)

	app.ServeHTTP(rec, req)

	// Try to get private without login

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/v1.0/private/whoami/", reader)

	app.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
