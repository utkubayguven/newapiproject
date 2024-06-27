package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DeleteUser godoc
// @Summary Delete a user
// @Scheme http
// @Tags User
// @Produce json
// @Accept json
// @Param id path int true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{id} [delete]
func (h Handler) DeleteUser(c *gin.Context) {
	userIDParam := c.Param("id") // Parametrelerden kullanıcı ID'sini al

	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı ID'si boş"})
		return
	}

	id, err := strconv.Atoi(userIDParam)
	if err != nil {
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

	// Kullanıcı kimliği eşleşmezse, işlemi iptal et
	if strconv.Itoa(id) != userIDString {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kullanıcıyı silme izniniz yok"})
		return
	}

	// etcd client'i al
	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	// etcd'den kullanıcıyı kontrol et
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "users/"+userIDString)
	if err != nil || resp.Count == 0 {
		fmt.Println("User not found or error retrieving user data:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	err = json.Unmarshal(resp.Kvs[0].Value, &user)
	if err != nil {
		fmt.Println("Error unmarshaling user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal user data: " + err.Error()})
		return
	}

	// etcd'den kullanıcıyı sil
	_, err = client.Delete(ctx, "users/"+userIDString)
	if err != nil {
		fmt.Println("Error deleting user data from etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user data from etcd: " + err.Error()})
		return
	}

	// Kullanıcının hesap bilgilerini de sil
	_, err = client.Delete(ctx, "accounts/"+userIDString)
	if err != nil {
		fmt.Println("Error deleting account data from etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete account data from etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil) // Her şey başarılı ise içerik olmadığını belirt
}
