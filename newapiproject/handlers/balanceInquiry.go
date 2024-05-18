package handlers

import (
	"net/http"
	"newapiprojet/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAccountBalance godoc
// @Summary Get the account balance
// @Description Get the account balance
// @Tags Account
// @Accept json
// @Produce json
// @Param accountNumber path int true "Account Number"
// @Success 200 {string} string "Balance inquiry successful"
// @Failure 404 {string} string "Account not found"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Router /balance/{accountNumber} [get]
func (h Handler) GetAccountBalance(c *gin.Context) {
	accountNumber := c.Param("accountNumber") // Extract account number from parameters
	var account models.Account

	// Query the account using AccountNumber
	id, err := strconv.Atoi(accountNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account number"})
		return
	}

	// JWT'den user_id'yi al
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme hatası"})
		return
	}

	// userID'yi uint olarak kontrol et
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
		return
	}

	// Check if the account belongs to the authenticated user
	if err := h.db.Where("id = ? AND user_id = ?", id, userIDUint).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found or access denied"})
		return
	}

	h.db.Create(&models.BalanceInquiry{AccountID: account.ID, CurrentBalance: account.Balance}) // Record the balance inquiry
	h.db.Save(&account)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Balance inquiry successful",
		"accountNumber": account.ID,
		"balance":       account.Balance,
	})
}
