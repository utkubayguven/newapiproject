package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const Bearer_schema = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			fmt.Println("Authorization header is missing.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authenticated!"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, Bearer_schema)
		if tokenString == authHeader {
			fmt.Println("Bearer token format is incorrect.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("Unexpected signing method:", token.Header["alg"])
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("utku123")), nil
		})

		if err != nil {
			fmt.Println("Token parsing or validation failed:", err)
			c.JSON(http.StatusForbidden, gin.H{"error": "Token is not valid!"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			fmt.Println("Token claims are invalid or token is expired.")
			c.JSON(http.StatusForbidden, gin.H{"error": "Token is not valid!"})
			c.Abort()
		}
	}
}

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
