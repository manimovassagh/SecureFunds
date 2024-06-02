package handlers

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"man-bank/models"
	"net/http"
)

type WithdrawHandler struct {
	DB *gorm.DB
}

type WithdrawRequest struct {
	AccountID uint    `json:"account_id" validate:"required"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
}

func (h *WithdrawHandler) Withdraw(c echo.Context) error {
	var request WithdrawRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	if err := c.Validate(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	var account models.Account
	if err := h.DB.First(&account, request.AccountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Account not found"})
	}

	if account.Balance < request.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Insufficient funds"})
	}

	account.Balance -= request.Amount
	if err := h.DB.Save(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update account balance"})
	}

	transaction := models.Transaction{
		AccountID: request.AccountID,
		Amount:    request.Amount,
		Type:      "withdrawal",
	}
	if err := h.DB.Create(&transaction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create transaction"})
	}

	return c.JSON(http.StatusOK, account)
}
