package notification

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	notificationGroup := router.Group("/api/notification")
	{
		notificationGroup.POST("", CreateNotification)
		notificationGroup.GET("/:id", ReadNotification)
		notificationGroup.GET("", ReadAllNotifications)
		notificationGroup.PUT("/:id", UpdateNotification)
		notificationGroup.DELETE("/:id", DeleteNotification)
	}
}
