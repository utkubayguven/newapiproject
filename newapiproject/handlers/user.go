package handlers

import (
	"net/http"
	"newapiprojet/models"
	"strconv"

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
	var user models.User
	userIDParam := c.Param("id") // Parametrelerden kullanıcı ID'sini al

	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, "Kullanıcı ID'si boş")
		return
	}

	id, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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

	// Kullanıcı kimliği eşleşmezse, işlemi iptal et
	if uint(id) != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kullanıcıyı silme izniniz yok"})
		return
	}

	// Kullanıcıyı sil
	if err := h.db.Where("id = ?", id).Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil) // Her şey başarılı ise içerik olmadığını belirt
}
