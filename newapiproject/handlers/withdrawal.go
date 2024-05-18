package handlers

import (
	"net/http"
	"newapiprojet/models"

	"github.com/gin-gonic/gin"
)

// Withdrawal godoc
// @Summary Withdraw money from an account
// @Description Withdraw money from an account
// @Tags Account
// @Accept json
// @Produce json
// @Param accountID path int true "Account ID"
// @Param withdrawalAmount path int true "Withdrawal Amount"
// @Success 200 {string} string "Withdrawal successful"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Account not found"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal Server Error"
// @Router /withdrawal [post]
func (h Handler) Withdrawal(c *gin.Context) {
	var account models.Account
	var withdrawal models.Withdrawal
	var input struct {
		AccountID        uint `json:"accountID"`
		WithdrawalAmount int  `json:"withdrawalAmount"`
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

	// userID'yi uint olarak kontrol et
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
		return
	}

	// Hesap bilgilerini al ve userID'yi kontrol et
	if err := h.db.Where("id = ?", input.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hesap bulunamadı"})
		return
	}

	if account.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
		return
	}

	if account.Balance < input.WithdrawalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetersiz bakiye"})
		return
	}

	account.Balance -= input.WithdrawalAmount
	withdrawal = models.Withdrawal{
		AccountID:        input.AccountID,
		WithdrawalAmount: input.WithdrawalAmount,
	}
	h.db.Create(&withdrawal)
	h.db.Save(&account)

	c.JSON(http.StatusOK, gin.H{
		"message": "Para çekme işlemi başarılı",
		"balance": account.Balance,
	})
}
