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

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, first name, last name, phone number, and PIN
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Success 201 {object} gin.H "User and Account created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 409 {string} string "Username already exists"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/register [post]
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
		fmt.Println("Error marshaling user data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal user data: " + err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.db.Get(ctx, "users/"+user.Username)
	if err != nil {
		fmt.Println("Error checking for existing username:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to check username: " + err.Error()})
		return
	}
	if resp != nil {
		fmt.Println("Username already exists:", user.Username)
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	err = h.db.Put(context.Background(), "users/"+user.Username, userData)
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

	accountData, err := json.Marshal(account)
	if err != nil {
		fmt.Println("Error marshaling account data:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal account data: " + err.Error()})
		return
	}

	err = h.db.Put(context.Background(), "accounts/"+user.Username, accountData)
	if err != nil {
		fmt.Println("Error storing account data in etcd:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store account data in etcd: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"account": account,
	})

	fmt.Println(user.ID, accountID)
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// Login godoc
// @Summary Login user and generate token
// @Description Login user and generate token
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body models.User true "User credentials"
// @Success 200 {string} string "Token generated"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var credentials models.User

	if err := c.BindJSON(&credentials); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.db.Get(ctx, "users/"+credentials.Username)
	if err != nil || resp == nil {
		fmt.Println("Invalid credentials or error retrieving user data:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	var user models.User
	err = json.Unmarshal(resp, &user)
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
