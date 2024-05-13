package handlers

import (
	"net/http"
	"newapiprojet/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /deleteuser/{id} [delete]
func (h Handler) DeleteUser(c *gin.Context) {
	var User models.User
	userId := c.Param("id") // Parametrelerden kullanıcı ID'sini al

	if userId == "" {
		c.JSON(http.StatusBadRequest, "Kullanıcı ID'si boş")
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.db.Where("id = ?", id).Delete(&User).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusNoContent, nil) // Her şey başarılı ise içerik olmadığını belirt
}
