package main

import (
	"testing"

	"newapiprojet/database"
	"newapiprojet/models"

	"gorm.io/gorm"
)

func cleanDatabase(db *gorm.DB) {
	// Drop tables in correct order due to foreign key constraints
	db.Exec("DROP TABLE IF EXISTS accounts")
	db.Exec("DROP TABLE IF EXISTS users")
}

func TestRegister(t *testing.T) {
	// Open a connection to the database
	db, err := database.InitDb()
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	// Clean the database before running the test
	cleanDatabase(db)

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Account{})

	// Create a new user
	user := models.User{
		Username:    "utku123",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "12345917890",
		PIN:         "1234",
	}

	// Create the user in the database
	result := db.Create(&user)
	if result.Error != nil {
		t.Fatalf("failed to create user: %v", result.Error)
	}

	// Check if the user was created successfully and ID is set
	if user.ID == 0 {
		t.Fatalf("expected user ID to be set, got 0")
	}

	// Create an account for the newly created user
	account := models.Account{
		UserID:  user.ID,
		Balance: 1000,
	}

	// Create the account in the database
	result = db.Create(&account)
	if result.Error != nil {
		t.Fatalf("failed to create account: %v", result.Error)
	}

	// Read the user
	var readUser models.User
	db.First(&readUser, "username = ?", user.Username)
	if readUser.Username != user.Username {
		t.Errorf("expected username to be %v, got %v", user.Username, readUser.Username)
	}

	// Read the account
	var readAccount models.Account
	db.First(&readAccount, "user_id = ?", user.ID)
	if readAccount.UserID != user.ID {
		t.Errorf("expected user_id to be %v, got %v", user.ID, readAccount.UserID)
	}

	// Update the user's PIN
	newPIN := "5678"
	db.Model(&readUser).Update("PIN", newPIN)
	db.First(&readUser, "username = ?", user.Username)
	if readUser.PIN != newPIN {
		t.Errorf("expected PIN to be %v, got %v", newPIN, readUser.PIN)
	}

	// Delete the account
	if err := db.Delete(&readAccount).Error; err != nil {
		t.Fatalf("failed to delete account: %v", err)
	}
	var countAccount int64
	db.Model(&models.Account{}).Where("user_id = ?", user.ID).Count(&countAccount)
	if countAccount != 0 {
		t.Errorf("expected account to be deleted, but count is %d", countAccount)
	} else {
		t.Logf("account successfully deleted")
	}

	// Delete the user
	if err := db.Delete(&readUser).Error; err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}
	var count int64
	db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count != 0 {
		t.Errorf("expected user to be deleted, but count is %d", count)
	} else {
		t.Logf("user successfully deleted")
	}

	// Clean the database after running the test
	cleanDatabase(db)
}

////////////////////////////////////////////
