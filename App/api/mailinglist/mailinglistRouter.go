package mailinglist

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	mailingListGroup := router.Group("/api/mailinglist")
	{
		// CRUD routes
		mailingListGroup.POST("", CreateMailinglist)
		mailingListGroup.GET("/:id", ReadMailinglist)
		mailingListGroup.GET("", ReadAllMailinglists)
		mailingListGroup.PUT("/:id", UpdateMailinglist)
		mailingListGroup.DELETE("/:id", DeleteMailinglist)

		// future routes
		mailingListGroup.POST("/:id/contact/:contact_id", AssignContactToMailingList)
		mailingListGroup.POST("/:id/tag/:tag_id", AssignTagToMailingList)

	}
}
