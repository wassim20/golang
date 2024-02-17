package contact

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	contactGroup := router.Group("/api/contact")
	{
		// CRUD routes
		contactGroup.POST("", CreateContact)
		contactGroup.GET("/:id", ReadContact)
		contactGroup.GET("", ReadAllContacts)
		contactGroup.PUT("/:id", UpdateContact)
		contactGroup.DELETE("/:id", DeleteContact)

		// future routes
		contactGroup.POST("/:id/tag/:tag_id", AssignTagToContact)

	}
}
