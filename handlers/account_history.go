package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"man-bank/models"
	"net/http"
)

type AccountHistoryHandler struct {
	DB *gorm.DB
}

func (h *AccountHistoryHandler) GetAccountHistory(c echo.Context) error {
	accountID := c.Param("account_id")

	var account models.Account
	if err := h.DB.Preload("Transactions").Preload("User").First(&account, accountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
	}

	response := map[string]interface{}{
		"username":     account.User.Username,
		"balance":      account.Balance,
		"transactions": account.Transactions,
	}

	return c.JSON(http.StatusOK, response)
}
