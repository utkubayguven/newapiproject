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

// Withdrawal godoc
// @Summary Withdraw money from an account
// @Description Withdraw money from an account
// @Tags Account
// @Accept json
// @Produce json
// @Param accountID path string true "Account ID"
// @Param withdrawalAmount body int true "Withdrawal Amount"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Account not found"
// @Failure 403 {string} string "Forbidden"
// @Router /withdrawal/{accountID} [post]
func (h Handler) Withdrawal(c *gin.Context) {
	var input struct {
		WithdrawalAmount int `json:"withdrawalAmount"`
	}
	accountID := c.Param("accountID") // Parametrelerden hesap ID'sini al

	// UUID olarak dönüştür
	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// JWT'den user_id'yi al
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme hatası"})
		return
	}

	// userID'yi string olarak kontrol et
	userIDString, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
		return
	}

	// etcd client'i al
	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	// etcd'den hesap bilgilerini al
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "accounts/"+accountUUID.String())
	if err != nil || resp.Count == 0 {
		fmt.Println("Account not found or error retrieving account data:", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	var account models.Account
	err = json.Unmarshal(resp.Kvs[0].Value, &account)
	if err != nil {
		fmt.Println("Error unmarshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal account data: " + err.Error()})
		return
	}

	// Hesap bilgilerini al ve userID'yi kontrol et
	if account.UserID.String() != userIDString {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
		return
	}

	if account.Balance < input.WithdrawalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetersiz bakiye"})
		return
	}

	account.Balance -= input.WithdrawalAmount
	withdrawal := models.Withdrawal{
		ID:               uuid.New(),
		AccountID:        account.ID,
		WithdrawalAmount: input.WithdrawalAmount,
		WithdrawalDate:   time.Now(),
	}
	withdrawalData, err := json.Marshal(withdrawal)
	if err != nil {
		fmt.Println("Error marshaling withdrawal data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal withdrawal data: " + err.Error()})
		return
	}

	_, err = client.Put(ctx, "withdrawals/"+withdrawal.ID.String(), string(withdrawalData))
	if err != nil {
		fmt.Println("Error storing withdrawal data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store withdrawal data in etcd: " + err.Error()})
		return
	}

	// Hesap bilgilerini güncelle
	accountData, err := json.Marshal(account)
	if err != nil {
		fmt.Println("Error marshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal account data: " + err.Error()})
		return
	}

	_, err = client.Put(ctx, "accounts/"+account.ID.String(), string(accountData))
	if err != nil {
		fmt.Println("Error storing account data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store account data in etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Para çekme işlemi başarılı",
		"balance": account.Balance,
	})
}
