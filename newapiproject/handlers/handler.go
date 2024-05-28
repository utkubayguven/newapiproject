package handlers

import (
	"net/http"
	"newapiprojet/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	h := Handler{db}
	return &h
}

func (h *Handler) SomeAPIHandler(c *gin.Context) {
	conf := config.GetConfig()
	err := conf.DecreaseRequestCount()
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Request limit exceeded"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request successful"})
}
