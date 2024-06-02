package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"man-bank/models"
)

type AccountHistoryHandler struct {
	DB *gorm.DB
}

func (h *AccountHistoryHandler) GetAccountHistory(c echo.Context) error {
	accountID := c.Param("account_id")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	var account models.Account
	if err := h.DB.Preload("Transactions").Preload("User").First(&account, accountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
	}

	if account.UserID != uint(userID) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
	}

	response := map[string]interface{}{
		"username":     account.User.Username,
		"balance":      account.Balance,
		"transactions": account.Transactions,
	}

	return c.JSON(http.StatusOK, response)
}
