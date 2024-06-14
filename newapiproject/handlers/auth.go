package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// // etcd client
// var etcdClient *clientv3.Client// get clıent tarzı fonksıyon olustur cunku timeouta dusuyor

// // InitEtcd initializes the etcd client
// func InitEtcd() {
// 	var err error
// 	etcdClient, err = clientv3.New(clientv3.Config{
// 		Endpoints:   []string{"http://etcd1:2379", "http://etcd2:2378", "http://etcd3:2377"},
// 		DialTimeout: 5 * time.Second,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var user models.User

	// Bind JSON from the request body to the user struct
	if err := c.BindJSON(&user); err != nil {
		fmt.Println("Error binding JSON:", err) // Log the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: Print the user data
	fmt.Printf("Received user data: %+v\n", user)

	// Trim any leading/trailing whitespace from the phone number
	user.PhoneNumber = strings.TrimSpace(user.PhoneNumber)

	// Validate the phone number has exactly 11 digits
	matchPhone, _ := regexp.MatchString(`^\d{11}$`, user.PhoneNumber)
	if !matchPhone {
		fmt.Println("Invalid phone number:", user.PhoneNumber) // Log the invalid phone number
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone must be exactly 11 digits"})
		return
	}

	// Validate the PIN has exactly 4 digits
	matchPin, _ := regexp.MatchString(`^\d{4}$`, user.PIN)
	if !matchPin {
		fmt.Println("Invalid PIN:", user.PIN) // Log the invalid PIN
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 4 digits"})
		return
	}

	// Generate a new UUID for the user
	user.ID = uuid.New()

	// Convert user struct to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error marshaling user data:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal user data: " + err.Error()})
		return
	}

	// Store user data in etcd
	_, err = h.client.Put(context.Background(), "users/"+user.Username, string(userData))
	if err != nil {
		fmt.Println("Error storing user data in etcd:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store user data in etcd: " + err.Error()})
		return
	}

	// Create an account for the newly registered user
	account := models.Account{
		UserID:  user.ID, // Set UserID to the newly created user's ID
		Balance: 1000,    // Default balance
	}

	// Convert account struct to JSON
	accountData, err := json.Marshal(account)
	if err != nil {
		fmt.Println("Error marshaling account data:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshal account data: " + err.Error()})
		return
	}

	// Store account data in etcd
	_, err = h.client.Put(context.Background(), "accounts/"+user.Username, string(accountData))
	if err != nil {
		fmt.Println("Error storing account data in etcd:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to store account data in etcd: " + err.Error()})
		return
	}

	// Return the newly created user and account information
	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"account": account,
	})
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}

// Login godoc
// @Summary Login user and generate token
// @Description Login user and generate token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.User true "User credentials"
// @Success 200 {string} string "Token generated"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
// func (h Handler) Login(c *gin.Context) {
// 	var credentials models.User

// 	// Bind JSON from request body
// 	if err := c.BindJSON(&credentials); err != nil {
// 		fmt.Println("Error binding JSON:", err) // Log the error
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get user data from etcd
// 	resp, err := etcdClient.Get(context.Background(), "users/"+credentials.Username)
// 	if err != nil || resp.Count == 0 {
// 		fmt.Println("Invalid credentials or error retrieving user data:", err) // Log the error
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	var user models.User
// 	err = json.Unmarshal(resp.Kvs[0].Value, &user)
// 	if err != nil {
// 		fmt.Println("Error unmarshaling user data:", err) // Log the error
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unmarshal user data: " + err.Error()})
// 		return
// 	}

// 	// Check if the PIN is correct
// 	if user.PIN != credentials.PIN {
// 		fmt.Println("Invalid PIN:", credentials.PIN) // Log the invalid PIN
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	// Create the JWT token
// 	expirationTime := time.Now().Add(7 * 24 * time.Hour) // Token expires in 7 days
// 	claims := &Claims{
// 		UserID: user.ID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 			Issuer:    "example.com",
// 			Subject:   credentials.Username,
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	secretKey := os.Getenv("JWT_SECRET")
// 	if secretKey == "" {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT secret key not configured"})
// 		return
// 	}
// 	tokenString, err := token.SignedString([]byte(secretKey))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
// 		return
// 	}

// 	// Return the token
// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }
