package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"man-bank/models"
)

var db *gorm.DB

func setup() {
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&models.User{}, &models.Account{}, &models.Transaction{})
	db = database

	truncateTables()
}

func truncateTables() {
	db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE accounts RESTART IDENTITY CASCADE")
	db.Exec("TRUNCATE TABLE transactions RESTART IDENTITY CASCADE")
}

func TestSignup(t *testing.T) {
	setup()
	e := echo.New()
	authHandler := &AuthHandler{DB: db}

	e.POST("/signup", authHandler.Signup)

	user := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonUser, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, authHandler.Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	setup()
	e := echo.New()
	authHandler := &AuthHandler{DB: db}

	// Ensure the user is created before login test
	db.Create(&models.User{Username: "testuser", Password: "$2a$10$bCLRQ5OGAq6ZIgxNXzm8YOx7hRVU4mbft6ldECQHLIz03fZcQcFOq"}) // bcrypt hash for "testpassword"

	e.POST("/login", authHandler.Login)

	user := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonUser, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, authHandler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestProtectedRoute(t *testing.T) {
	setup()
	e := echo.New()
	authHandler := &AuthHandler{DB: db}

	// Ensure the user is created before login test
	db.Create(&models.User{Username: "testuser", Password: "$2a$10$bCLRQ5OGAq6ZIgxNXzm8YOx7hRVU4mbft6ldECQHLIz03fZcQcFOq"}) // bcrypt hash for "testpassword"

	e.POST("/login", authHandler.Login)

	user := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonUser, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, authHandler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		json.Unmarshal(rec.Body.Bytes(), &response)
		token := response["token"]

		req = httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		c.SetPath("/protected")

		jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: []byte(os.Getenv("JWT_SECRET")),
		})

		handler := func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			userId, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid user ID"})
			}
			username, ok := claims["username"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username"})
			}

			return c.JSON(http.StatusOK, map[string]interface{}{
				"message":  "You are in a protected route!",
				"user_id":  userId,
				"username": username,
			})
		}

		h := jwtMiddleware(handler)

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	}
}
