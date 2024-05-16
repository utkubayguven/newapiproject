package handlers

import (
	"net/http"
	"newapiprojet/models"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Success 201 "User created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h Handler) Register(c *gin.Context) {
	var user models.User

	// Bind JSON from the request body to the user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the phone number has exactly 11 digits
	matchPhone, _ := regexp.MatchString(`^\d{11}$`, user.PhoneNumber)
	if !matchPhone {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone must be exactly 11 digits"})
		return
	}

	// Validate the PIN has exactly 4 digits
	matchPin, _ := regexp.MatchString(`^\d{4}$`, user.PIN)
	if !matchPin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 4 digits"})
		return
	}

	// Create the user in the database
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create user: " + err.Error()})
		return
	}

	// Create an account for the newly registered user
	account := models.Account{
		UserID:  user.ID, // Set UserID to the newly created user's ID
		Balance: 1000,    // Default balance
	}

	// Create the account in the database
	if err := h.db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create account: " + err.Error()})
		return
	}

	// Return the newly created user and account information
	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"account": account,
	})
}

// Login godoc
// @Summary Login to the application
// @Description Login to the application
// @Tags User
// @Accept json
// @Produce json
// @Param username body string true "Username"
// @Param pin body string true "PIN"
// @Success 200 {string} string "Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (h Handler) Login(c *gin.Context) {
	var credentials models.User // Assuming you have a Credentials model with fields `Username` and `Password`

	// Bind JSON from request body
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.User
	if err := h.db.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check if the password is correct
	if user.PIN != credentials.PIN {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Create the JWT token
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // Token expires in 7 days
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "example.com",
		Subject:   credentials.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT gizli anahtarı yapılandırılmamış"})
		return
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
