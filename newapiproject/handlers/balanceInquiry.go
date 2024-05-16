package handlers

import (
	"net/http"
	"newapiprojet/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func (h Handler) BalanceInquiry(c *gin.Context) {
// 	userId, _ := c.Get("userID")
// 	var account models.Account

// 	if err := h.db.Where("user_id = ?", userId.(uint)).First(&account).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
// 		return
// 	}

// 	newBalanceInquiry := models.BalanceInquiry{
// 		AccountID:      account.ID,
// 		CurrentBalance: account.Balance,
// 	}

// 	h.db.Create(&newBalanceInquiry) // Record the balance inquiry

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":        "Balance inquiry successful",
// 		"accountNumber":  account.ID,
// 		"currentBalance": account.Balance,
// 	})
// }

// GetAccountBalance godoc
// @Summary Get the account balance
// @Description Get the account balance
// @Tags Account
// @Accept json
// @Produce json
// @Param accountNumber path int true "Account Number"
// @Success 200 {string} string "Balance inquiry successful"
// @Failure 404 {string} string "Account not found"
// @Failure 400 {string} string "Bad Request"
// @Router /balance/{accountNumber} [get]
func (h Handler) GetAccountBalance(c *gin.Context) {
	accountNumber := c.Param("accountNumber") // Extract account number from parameters
	var account models.Account

	// Query the account using AccountNumber
	id, err := strconv.Atoi(accountNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account number"})
		return
	}

	if err := h.db.Where("id = ?", id).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Account not found"})
		return
	}

	h.db.Create(&models.BalanceInquiry{AccountID: account.ID, CurrentBalance: account.Balance}) // Record the balance inquiry
	h.db.Save(&account)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Balance inquiry successful",
		"accountNumber": account.ID,
		"balance":       account.Balance,
	})

}
