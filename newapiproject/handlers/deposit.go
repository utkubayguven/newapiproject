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

// Deposit godoc
// @Summary Deposit money into an account
// @Description Deposit money into an account
// @Tags Account
// @Accept json
// @Produce json
// @Param accountID path string true "Account ID"
// @Param depositAmount body int true "Deposit Amount"
// @Success 200 {string} string "Deposit successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Account not found"
// @Failure 403 {string} string "Forbidden"
// @Router /deposit/{accountID} [post]
func (h Handler) Deposit(c *gin.Context) {
	var input struct {
		DepositAmount int `json:"depositAmount"`
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

	if input.DepositAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit amount must be positive"})
		return
	}

	account.Balance += input.DepositAmount
	deposit := models.Deposit{
		ID:            uuid.New(),
		AccountID:     account.ID,
		DepositAmount: input.DepositAmount,
		DepositDate:   time.Now(),
	}
	depositData, err := json.Marshal(deposit)
	if err != nil {
		fmt.Println("Error marshaling deposit data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal deposit data: " + err.Error()})
		return
	}

	_, err = client.Put(ctx, "deposits/"+deposit.ID.String(), string(depositData))
	if err != nil {
		fmt.Println("Error storing deposit data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store deposit data in etcd: " + err.Error()})
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
		"message": "Deposit successful",
		"balance": account.Balance,
	})
}
