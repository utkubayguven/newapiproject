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

// GetAccountByID godoc
// @Summary Get an account by ID
// @Description Get an account by ID
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} models.Account "Account found"
// @Failure 404 {string} string "Account not found"
// @Failure 400 {string} string "Bad Request"
// @Router /account/{id} [get]
func (h Handler) GetAccountByID(c *gin.Context) {
	accountID := c.Param("id")
	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID format"})
		return
	}

	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "accounts/"+accountUUID.String())
	if err != nil || resp.Count == 0 {
		fmt.Println("Account not found or error retrieving account data:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var account models.Account
	err = json.Unmarshal(resp.Kvs[0].Value, &account)
	if err != nil {
		fmt.Println("Error unmarshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal account data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// DeleteAccountByID godoc
// @Summary Delete an account by ID
// @Description Delete an account by ID
// @Tags Account
// @Produce json
// @Param id path string true "Account ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Account not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /account/{id} [delete]
func (h Handler) DeleteAccountByID(c *gin.Context) {
	accountID := c.Param("id")

	if accountID == "" {
		c.JSON(http.StatusBadRequest, "Account ID cannot be empty")
		return
	}

	// Get etcd client
	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// JWT'den user_id'yi al
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication error"})
		return
	}

	// userID'yi uuid.UUID olarak kontrol et
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID conversion error"})
		return
	}

	resp, err := client.Get(ctx, "accounts/"+accountID)
	if err != nil {
		fmt.Println("Error retrieving account data:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if resp.Count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var account models.Account
	err = json.Unmarshal(resp.Kvs[0].Value, &account)
	if err != nil {
		fmt.Println("Error unmarshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal account data: " + err.Error()})
		return
	}

	// Hesabın kullanıcıya ait olup olmadığını kontrol et
	if account.UserID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this account"})
		return
	}

	// Hesabı sil
	_, err = client.Delete(ctx, "accounts/"+accountID)
	if err != nil {
		fmt.Println("Error deleting account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete account data: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
