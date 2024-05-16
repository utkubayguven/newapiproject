package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (m Newapiprojetmiddlewares) LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		bodysize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		if errorMessage != "" {
			errorMessage = fmt.Sprintf("Error: %s", errorMessage)
		}

		if statusCode >= 400 {
			fmt.Printf("[ERROR][%s] [%s] [%d] [%s] [%d] [%s] [%d] [%s]\n", method, path, statusCode, clientIP, bodysize, latency, statusCode, errorMessage)
		} else {
			fmt.Printf("[INFO][%s] [%s] [%d] [%s] [%d] [%s] [%d] [%s]\n", method, path, statusCode, clientIP, bodysize, latency, statusCode, errorMessage)
		}
	}
}
