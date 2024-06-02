package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"man-bank/models"
)

type TransferHandler struct {
	DB *gorm.DB
}

type TransferRequest struct {
	FromAccountID uint    `json:"from_account_id" validate:"required"`
	ToAccountID   uint    `json:"to_account_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

func (h *TransferHandler) Transfer(c echo.Context) error {
	var request TransferRequest
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

	var fromAccount, toAccount models.Account
	if err := h.DB.First(&fromAccount, request.FromAccountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Sender account not found"})
	}
	if err := h.DB.First(&toAccount, request.ToAccountID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Receiver account not found"})
	}

	if fromAccount.UserID != uint(userID) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
	}

	if fromAccount.Balance < request.Amount {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Insufficient funds"})
	}

	// Perform the transfer within a transaction
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		fromAccount.Balance -= request.Amount
		if err := tx.Save(&fromAccount).Error; err != nil {
			return err
		}

		toAccount.Balance += request.Amount
		if err := tx.Save(&toAccount).Error; err != nil {
			return err
		}

		// Log transactions for both accounts
		fromTransaction := models.Transaction{
			AccountID: fromAccount.ID,
			Amount:    request.Amount,
			Type:      "transfer_out",
		}
		if err := tx.Create(&fromTransaction).Error; err != nil {
			return err
		}

		toTransaction := models.Transaction{
			AccountID: toAccount.ID,
			Amount:    request.Amount,
			Type:      "transfer_in",
		}
		if err := tx.Create(&toTransaction).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Transfer failed"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"from_account": fromAccount,
		"to_account":   toAccount,
	})
}
