package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAccountByID godoc
// @Summary Get an account by ID
// @Description Get an account by ID
// @Tags Account
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {string} string "Account found"
// @Failure 404 {string} string "Account not found"
// @Failure 400 {string} string "Bad Request"
// @Router /account/{id} [get]
func (h Handler) GetAccountByID(c *gin.Context) {
	var account models.Account
	accountID := c.Param("id")

	id, err := strconv.Atoi(accountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz hesap ID formatı"})
		return
	}

	fmt.Println("ID:", id)

	if err := h.db.Preload("Deposits").Preload("Withdrawals").Preload("BalanceInquiries").First(&account, "id = ?", id).Error; err != nil {
		fmt.Println("DB Error:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Hesap bulunamadı")
			c.JSON(http.StatusNotFound, gin.H{"error": "Hesap bulunamadı"})
		} else {
			fmt.Println("Başka bir hata oluştu:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Veri tabanı hatası"})
		}
		return
	}

	fmt.Println("Account found:", account)

	c.JSON(http.StatusOK, gin.H{
		"account": account,
	})
}

// DeleteAccount godoc
// @Summary Delete an account
// @Scheme http
// @Tags Account
// @Produce json
// @Accept json
// @Param accountNumber path int true "Account Number"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Account not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /deleteacc/{accountNumber} [delete]
func (h Handler) DeleteAccount(c *gin.Context) {
	var account models.Account
	accountNumber := c.Param("accountNumber")

	if accountNumber == "" {
		c.JSON(http.StatusBadRequest, "Hesap numarası boş")
		return
	}

	id, err := strconv.Atoi(accountNumber)
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

	// Hesabın kullanıcıya ait olup olmadığını kontrol et
	if err := h.db.Where("id = ? AND user_id = ?", id, userIDUint).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found or access denied"})
		return
	}

	// Hesabı sil
	if err := h.db.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil) // Her şey başarılı ise içerik olmadığını belirt
}
