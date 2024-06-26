package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register handler
func (h *Handler) Register(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received user data: %+v\n", user)

	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)

	matchPhone, _ := regexp.MatchString(`^\d{11}$`, user.PhoneNumber)
	if !matchPhone {
		fmt.Println("Invalid phone number:", user.PhoneNumber)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone must be exactly 11 digits"})
		return
	}

	matchPin, _ := regexp.MatchString(`^\d{4}$`, user.PIN)
	if !matchPin {
		fmt.Println("Invalid PIN:", user.PIN)
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 4 digits"})
		return
	}

	user.ID = uuid.New()

	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user data:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal user data: " + err.Error()})
		return
	}

	// Get etcd client
	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "users/"+user.Username)
	if err != nil {
		fmt.Println("Error checking for existing username:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to check username: " + err.Error()})
		return
	}
	if resp.Count > 0 {
		fmt.Println("Username already exists:", user.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Store user data in etcd
	_, err = client.Put(context.Background(), "users/"+user.Username, string(userData))
	if err != nil {
		fmt.Println("Error storing user data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store user data in etcd: " + err.Error()})
		return
	}

	accountID := uuid.New()

	account := models.Account{
		ID:      accountID,
		UserID:  user.ID,
		Balance: 1000,
	}

	// Convert account struct to JSON
	accountData, err := json.Marshal(account)
	if err != nil {
		fmt.Println("Error marshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal account data: " + err.Error()})
		return
	}

	// Store account data in etcd
	_, err = client.Put(context.Background(), "accounts/"+user.Username, string(accountData))
	if err != nil {
		fmt.Println("Error storing account data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store account data in etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"account": account,
	})
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

func (h *Handler) Login(c *gin.Context) {
	var credentials models.User

	if err := c.BindJSON(&credentials); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "users/"+credentials.Username)
	if err != nil || resp.Count == 0 {
		fmt.Println("Invalid credentials or error retrieving user data:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var user models.User
	err = json.Unmarshal(resp.Kvs[0].Value, &user)
	if err != nil {
		fmt.Println("Error unmarshaling user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal user data: " + err.Error()})
		return
	}

	if user.PIN != credentials.PIN {
		fmt.Println("Invalid PIN:", credentials.PIN)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "example.com",
			Subject:   credentials.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret key not configured"})
		return
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
