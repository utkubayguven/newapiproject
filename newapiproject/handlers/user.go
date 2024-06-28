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

// DeleteUser godoc
// @Summary Delete a user
// @Scheme http
// @Tags User
// @Produce json
// @Accept json
// @Param id path string true "User ID"
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

	userUUID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
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

	fmt.Printf("userUUID: %s, userIDUUID: %s\n", userUUID.String(), userIDUUID.String())

	if userUUID != userIDUUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kullanıcıyı silme izniniz yok"})
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

	resp, err := client.Get(ctx, "users/"+userUUID.String())
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

	_, err = client.Delete(ctx, "users/"+userUUID.String())
	if err != nil {
		fmt.Println("Error deleting user data from etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user data from etcd: " + err.Error()})
		return
	}

	_, err = client.Delete(ctx, "accounts/"+userUUID.String())
	if err != nil {
		fmt.Println("Error deleting account data from etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete account data from etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil) //
}
