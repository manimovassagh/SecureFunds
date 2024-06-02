package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"man-bank/models"
)

type DepositHandler struct {
	DB *gorm.DB
}

type DepositRequest struct {
	AccountID uint    `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

func (h *DepositHandler) Deposit(c echo.Context) error {
	var request DepositRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid token"})
	}

	var account models.Account
	if err := h.DB.First(&account, request.AccountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
	}

	if account.UserID != uint(userID) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
	}

	account.Balance += request.Amount
	if err := h.DB.Save(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update account balance"})
	}

	transaction := models.Transaction{
		AccountID: request.AccountID,
		Amount:    request.Amount,
		Type:      "deposit",
	}
	if err := h.DB.Create(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create transaction"})
	}

	return c.JSON(http.StatusOK, account)
}
