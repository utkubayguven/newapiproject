package handlers

// import (
// 	"net/http"
// 	"newapiprojet/models"
// 	"regexp"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // PinChange godoc
// // @Summary Change the user's PIN
// // @Description Change the user's PIN
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param id path int true "User ID"
// // @Param input body struct{OldPIN string `json:"oldPIN";NewPIN string `json:"newPIN"`} true "Old and New PIN"
// // @Success 200 {string} string "PIN updated successfully"
// // @Failure 400 {string} string "Bad Request"
// // @Failure 404 {string} string "User not found"
// // @Failure 403 {string} string "Forbidden"
// // @Failure 500 {string} string "Internal Server Error"
// // @Router /pin-change/{id} [post]
// func (h Handler) PinChange(c *gin.Context) {
// 	var user models.User
// 	var input struct { // Define a struct to bind the input JSON
// 		OldPIN string `json:"oldPIN"`
// 		NewPIN string `json:"newPIN"`
// 	}
// 	idParam := c.Param("id")

// 	id, err := strconv.Atoi(idParam)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
// 		return
// 	}

// 	// JWT'den user_id'yi al
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Yetkilendirme hatası"})
// 		return
// 	}

// 	// userID'yi uint olarak kontrol et
// 	userIDUint, ok := userID.(uint)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID dönüştürme hatası"})
// 		return
// 	}

// 	// Kullanıcı kimliği eşleşmezse, işlemi iptal et
// 	if uint(id) != userIDUint {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "Bu hesaba erişim izniniz yok"})
// 		return
// 	}

// 	// Retrieve user from the database
// 	if err := h.db.First(&user, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Bind the JSON payload to the input struct
// 	if err := c.BindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Validate the NewPIN is exactly 4 digits
// 	matchPin, _ := regexp.MatchString(`^\d{4}$`, input.NewPIN)
// 	if !matchPin {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN must be exactly 4 digits"})
// 		return
// 	}

// 	// Check if the provided OldPIN matches the current PIN
// 	if input.OldPIN != user.PIN {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "OldPIN does not match the current PIN"})
// 		return
// 	}

// 	// Update the user's PIN in the database
// 	user.PIN = input.NewPIN
// 	if err := h.db.Save(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update PIN"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "PIN updated successfully"})
// }
