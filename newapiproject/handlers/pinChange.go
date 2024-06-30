package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PinChange godoc
// @Summary Change the user's PIN
// @Description Change the user's PIN
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param input body struct{OldPIN string `json:"oldPIN";NewPIN string `json:"newPIN"`} true "Old and New PIN"
// @Success 200 {string} string "PIN updated successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User not found"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /pin-change/{id} [post]
func (h *Handler) PinChange(c *gin.Context) {
	var input struct {
		OldPIN string `json:"oldPIN"`
		NewPIN string `json:"newPIN"`
	}
	idParam := c.Param("id")

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
	userIDString := userIDUUID.String()

	fmt.Printf("idParam: %s, userIDString: %s\n", idParam, userIDString)

	if idParam != userIDString {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := h.db.Get(ctx, "users/"+userIDString)
	if err != nil || data == nil {
		fmt.Println("User not found or error retrieving user data:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	err = json.Unmarshal(data, &user)
	if err != nil {
		fmt.Println("Error unmarshaling user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal user data: " + err.Error()})
		return
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	matchPin, _ := regexp.MatchString(`^\d{4}$`, input.NewPIN)
	if !matchPin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 4 digits"})
		return
	}

	if input.OldPIN != user.PIN {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OldPIN does not match the current PIN"})
		return
	}

	user.PIN = input.NewPIN
	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal user data: " + err.Error()})
		return
	}

	err = h.db.Put(ctx, "users/"+userIDString, userData)
	if err != nil {
		fmt.Println("Error updating user data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user data in etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PIN updated successfully"})
}
