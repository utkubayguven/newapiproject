package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAccountBalance godoc
// @Summary Get the account balance
// @Description Get the account balance
// @Tags Account
// @Accept json
// @Produce json
// @Param accountID path string true "Account ID"
// @Success 200 {string} string "Balance inquiry successful"
// @Failure 404 {string} string "Account not found"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Router /balance/{accountID} [get]
func (h Handler) GetAccountBalance(c *gin.Context) {
	accountID := c.Param("accountID")

	if accountID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account ID cannot be empty"})
		return
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme hatası"})
		return
	}

	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.db.Get(ctx, "accounts/username"+accountUUID.String())
	if err != nil || resp == nil {
		fmt.Println("Account not found or error retrieving account data:", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	var account models.Account
	err = json.Unmarshal(resp, &account)
	if err != nil {
		fmt.Println("Error unmarshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal account data: " + err.Error()})
		return
	}

	if account.UserID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
		return
	}

	balanceInquiry := models.BalanceInquiry{
		ID:             uuid.New(),
		AccountID:      account.ID,
		CurrentBalance: account.Balance,
		InquiryDate:    time.Now(),
	}
	inquiryData, err := json.Marshal(balanceInquiry)
	if err != nil {
		fmt.Println("Error marshaling balance inquiry data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal balance inquiry data: " + err.Error()})
		return
	}

	err = h.db.Put(ctx, "balance_inquiries/"+accountUUID.String(), inquiryData)
	if err != nil {
		fmt.Println("Error storing balance inquiry data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store balance inquiry data in etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Balance inquiry successful",
		"accountID": account.ID,
		"balance":   account.Balance,
	})
}
