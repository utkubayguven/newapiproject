package middlewares

import (
	"net/http"
	"newapiprojet/config"

	"github.com/gin-gonic/gin"
)

func RequestLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfig()
		err := conf.DecreaseRequestCount()
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Request limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
