package handlers

// import (
// 	"net/http"
// 	"newapiprojet/models"

// 	"github.com/gin-gonic/gin"
// )

// // Deposit godoc
// // @Summary Deposit money into an account
// // @Description Deposit money into an account
// // @Tags Account
// // @Accept json
// // @Produce json
// // @Param accountID path int true "Account ID"
// // @Param depositAmount path int true "Deposit Amount"
// // @Success 200 {string} string "Deposit successful"
// // @Failure 400 {string} string "Bad Request"
// // @Failure 404 {string} string "Account not found"
// // @Failure 500 {string} string "Internal Server Error"
// // @Router /deposit [post]
// func (h Handler) Deposit(c *gin.Context) {
// 	var account models.Account
// 	var deposit models.Deposit
// 	var input struct {
// 		AccountID     uint `json:"accountID"`
// 		DepositAmount int  `json:"depositAmount"`
// 	}

// 	if err := c.BindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// JWT'den user_id'yi al
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme hatası"})
// 		return
// 	}

// 	// userID'yi uint olarak kontrol et
// 	userIDUint, ok := userID.(uint)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
// 		return
// 	}

// 	// Hesap bilgilerini al ve userID'yi kontrol et
// 	if err := h.db.Where("id = ?", input.AccountID).First(&account).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Hesap bulunamadı"})
// 		return
// 	}

// 	if account.UserID != userIDUint {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
// 		return
// 	}

// 	if input.DepositAmount <= 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Deposit amount must be positive"})
// 		return
// 	}

// 	account.Balance += input.DepositAmount
// 	deposit = models.Deposit{
// 		AccountID:     input.AccountID,
// 		DepositAmount: input.DepositAmount,
// 	}

// 	h.db.Create(&deposit)
// 	h.db.Save(&account)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Deposit successful",
// 		"balance": account.Balance,
// 	})
// }
