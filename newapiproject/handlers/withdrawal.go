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

	if err := h.db.Where("id = ?", input.AccountID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if account.Balance < input.WithdrawalAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
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
		"message": "Withdrawal successful",
		"balance": account.Balance,
	})
}
