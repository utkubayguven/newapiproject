package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"newapiprojet/models"
	"regexp"
	"strings"
	"time"

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

	// Get etcd client
	client, err := h.getClient()
	if err != nil {
		fmt.Println("Error getting etcd client:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get etcd client: " + err.Error()})
		return
	}

	// Check if username already exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "users/"+user.Username)
	if err != nil {
		fmt.Println("Error checking for existing username:", err) // Log the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to check username: " + err.Error()})
		return
	}
	if resp.Count > 0 {
		fmt.Println("Username already exists:", user.Username) // Log the duplicate username
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
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
	_, err = client.Put(ctx, "users/"+user.Username, string(userData))
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
	_, err = client.Put(ctx, "accounts/"+user.Username, string(accountData))
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

// Claims struct for JWT
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.StandardClaims
}
