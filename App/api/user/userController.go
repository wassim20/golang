package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassim_p/App/initializers"
	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context) {
	var service = NewUserService(initializers.DB)
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Create(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func ReadUser(c *gin.Context) {
	var service = NewUserService(initializers.DB)

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read user: %v", err)})
		return
	}

	user, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read user: %v", err)})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func ReadAllUsers(c *gin.Context) {
	var service = NewUserService(initializers.DB)

	users, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read users: %v", err)})
		return
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context) {
	var service = NewUserService(initializers.DB)
	var user models.User

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read user: %v", err)})
		return
	}
	user.ID = uint(idUint)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Update(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to update user: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func DeleteUser(c *gin.Context) {
	var service = NewUserService(initializers.DB)
	var user models.User

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete user: %v", err)})
		return
	}
	user.ID = uint(idUint)

	if err := service.Delete(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete user: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, "User deleted successfully")
}

func AssignNotificationToUser(c *gin.Context) {
	var service = NewUserService(initializers.DB)

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse user ID: %v", err)})
		return
	}

	notificationID, err := strconv.ParseUint(c.Param("notification_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse notification ID: %v", err)})
		return
	}

	err = service.AssignNotification(uint(userID), uint(notificationID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign notification to user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification assigned to user successfully"})
}
