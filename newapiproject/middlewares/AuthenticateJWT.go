package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme başlığı eksik."})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, BearerSchema)
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token formatı"})
			c.Abort()
			return
		}

		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT gizli anahtarı yapılandırılmamış"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("beklenmedik imzalama yöntemi: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Token geçerli değil"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("userID", claims["user_id"])
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Token geçerli değil"})
			c.Abort()
		}
	}
}
