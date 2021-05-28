package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUsers(t *testing.T) {
	t.Run("Test create user invalid data", func(t *testing.T) {
		body := `{
			"username":"knz@email.com,
			"password":"1234"
			}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader((body)))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("Test Create user", func(t *testing.T) {
		// var user User
		body := `
			{
				"username":"knz@email.com",
				"password":"12345678"
			}	
		`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader((body)))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.Code)
		token := res.Header().Get("X-Auth-token")
		assert.NotEmpty(t, token)
		// err = json.Unmarshal(res.Body.Bytes(), &user)
		// t.Logf("user: %v", res.Body.String())
		// assert.Nil(t, err)
		// assert.Equal(t, "knz@email.com", user.Email)
		// assert.Empty(t, user.Password)
	})

	t.Run("Test Create duplicate", func(t *testing.T) {
		body := `
			{
				"username":"knz@email.com",
				"password":"12345678"
			}	
		`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader((body)))
		res := httptest.NewRecorder()
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		e := echo.New()
		c := e.NewContext(req, res)
		uh.Col = usersCol
		err := uh.CreateUser(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.Code)
	})
}
