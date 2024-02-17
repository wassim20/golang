package notification

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

func CreateNotification(c *gin.Context) {
	var service = NewNotificationService(initializers.DB)
	var notification models.Notification

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Create(&notification); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create notification: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func ReadNotification(c *gin.Context) {
	var service = NewNotificationService(initializers.DB)

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read notification: %v", err)})
		return
	}

	notification, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read notification: %v", err)})
		return
	}

	if notification == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, notification)
}

func ReadAllNotifications(c *gin.Context) {
	var service = NewNotificationService(initializers.DB)

	notifications, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read notifications: %v", err)})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func UpdateNotification(c *gin.Context) {
	var service = NewNotificationService(initializers.DB)
	var notification models.Notification

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read notification: %v", err)})
		return
	}
	notification.ID = uint(idUint)

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Update(&notification); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to update notification: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, notification)
}

func DeleteNotification(c *gin.Context) {
	var service = NewNotificationService(initializers.DB)
	var notification models.Notification

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete notification: %v", err)})
		return
	}
	notification.ID = uint(idUint)

	if err := service.Delete(&notification); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete notification: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, "Notification deleted successfully")
}
