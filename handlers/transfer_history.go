package handlers

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"man-bank/models"
)

type TransferHistoryHandler struct {
	DB *gorm.DB
}

func (h *TransferHistoryHandler) GetTransferHistory(c echo.Context) error {
	accountID := c.Param("account_id")

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	var account models.Account
	if err := h.DB.Preload("User").First(&account, accountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
	}

	if account.UserID != uint(userID) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
	}

	// Pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	var transactions []models.Transaction
	if err := h.DB.Where("account_id = ? AND type IN ?", accountID, []string{"transfer_in", "transfer_out"}).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&transactions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve transactions"})
	}

	response := map[string]interface{}{
		"username":     account.User.Username,
		"transactions": transactions,
		"page":         page,
		"limit":        limit,
	}

	return c.JSON(http.StatusOK, response)
}
