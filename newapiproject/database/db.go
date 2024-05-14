package database

import (
	"fmt"
	"newapiprojet/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDb initializes the database connection and runs migrations
func InitDb() (*gorm.DB, error) {
	dbURL := "postgres://root:root@localhost:5432/test_db"
	DB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	err = DB.AutoMigrate(&models.User{}, &models.Account{}, &models.Deposit{}, &models.Withdrawal{}, &models.BalanceInquiry{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return DB, nil
}
