package tag

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	tagGroup := router.Group("/api/tag")
	{
		// CRUD routes
		tagGroup.POST("", CreateTag)
		tagGroup.GET("/:id", ReadTag)
		tagGroup.GET("", ReadAllTags)
		tagGroup.PUT("/:id", UpdateTag)
		tagGroup.DELETE("/:id", DeleteTag)

		// Placeholder for future routes
		// tagGroup.POST("/:id/mailinglists", AssignMailingList)
		// tagGroup.POST("/:id/contacts", AssignContact)
	}
}
