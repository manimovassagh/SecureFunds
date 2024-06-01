package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"man-bank/models"
)

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Signup(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := user.HashPassword(user.Password); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error hashing password"})
	}

	if err := h.DB.Create(user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating user"})
	}

	// Automatically create an account for the new user
	account := models.Account{
		UserID:      user.ID,
		AccountType: "checking", // or any default account type
		Balance:     0,
	}

	if err := h.DB.Create(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating account"})
	}

	// Load user with the associated accounts
	if err := h.DB.Preload("Accounts").First(user, user.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error loading user with accounts"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	dbUser := new(models.User)
	if err := h.DB.Where("username = ?", user.Username).First(dbUser).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	if err := dbUser.CheckPassword(user.Password); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
	}

	token, err := createJWTToken(dbUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error generating token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func createJWTToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
