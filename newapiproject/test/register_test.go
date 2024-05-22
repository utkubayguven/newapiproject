package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"newapiprojet/handlers"
	"newapiprojet/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestRegister(t *testing.T) {
	// Initialize sqlmock
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlDB.Close()

	// Initialize GORM with the sqlmock database
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Günlüklemeyi etkinleştir
	})
	if err != nil {
		t.Fatalf("failed to initialize gorm db: %v", err)
	}

	handler := handlers.NewHandler(gormDB)

	// Create Gin router
	router := gin.Default()
	router.POST("/register", handler.Register)

	// Create a new user
	user := models.User{
		Username:    "utku123",
		FirstName:   "John",
		LastName:    "Doe",
		PhoneNumber: "12345917890",
		PIN:         "1234",
	}

	// Set expectations for the sqlmock
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","username","first_name","last_name","phone_number","pin"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) RETURNING "id"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, user.FirstName, user.LastName, user.PhoneNumber, user.PIN).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectExec(`INSERT INTO "accounts" \("created_at","updated_at","deleted_at","user_id","balance"\) VALUES \(\$1,\$2,\$3,\$4,\$5\)`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 1000).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Create a valid register request
	reqBody, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
