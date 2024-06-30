package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"newapiprojet/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Withdrawal godoc
// @Summary Withdraw money from an account
// @Description Withdraw money from an account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body struct{ AccountID uuid.UUID `json:"accountID"`; WithdrawalAmount int `json:"withdrawalAmount"` } true "Withdrawal details"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Account not found"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /account/withdrawal [post]
func (h *Handler) Withdrawal(c *gin.Context) {
	var input struct {
		AccountID        uuid.UUID `json:"accountID"`
		WithdrawalAmount int       `json:"withdrawalAmount"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication error"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID conversion error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.db.Get(ctx, "accounts/"+input.AccountID.String())
	if err != nil || resp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var account models.Account
	err = json.Unmarshal(resp, &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal account data"})
		return
	}

	if account.UserID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if account.Balance < input.WithdrawalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	account.Balance -= input.WithdrawalAmount

	accountData, err := json.Marshal(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal account data"})
		return
	}

	err = h.db.Put(context.Background(), "accounts/"+account.ID.String(), accountData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update account data"})
		return
	}

	withdrawal := models.Withdrawal{
		ID:               uuid.New(),
		AccountID:        account.ID,
		WithdrawalAmount: input.WithdrawalAmount,
		WithdrawalDate:   time.Now(),
	}

	withdrawalData, err := json.Marshal(withdrawal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal withdrawal data"})
		return
	}

	err = h.db.Put(context.Background(), "withdrawals/"+withdrawal.ID.String(), withdrawalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store withdrawal data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Withdrawal successful",
		"balance": account.Balance,
	})
}
