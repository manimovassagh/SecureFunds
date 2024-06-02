package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"man-bank/handlers"
	"man-bank/models"
)

var (
	db  *gorm.DB
	err error
	e   *echo.Echo
)

func init() {
	dsn := "host=localhost user=postgres password=postgres dbname=bank port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})
	e = echo.New()
}

func TestSignup(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Password: "testpassword",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.AuthHandler{DB: db}
	if assert.NoError(t, handler.Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, user.Username, response["username"])
	}
}

func TestLogin(t *testing.T) {
	user := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonUser, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := handlers.AuthHandler{DB: db}
	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
	}
}
