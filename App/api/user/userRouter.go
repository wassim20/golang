package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("", CreateUser)
		userGroup.GET("/:id", ReadUser)
		userGroup.GET("", ReadAllUsers)
		userGroup.PUT("/:id", UpdateUser)
		userGroup.DELETE("/:id", DeleteUser)
		userGroup.POST("/:id/notification/:notification_id", AssignNotificationToUser)

	}
}
