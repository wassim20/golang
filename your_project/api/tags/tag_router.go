package tags

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TagRouterInit initializes the routes related to tags.
func TagRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewTagRepository(db)

	// Private
	tag := router.Group("/:companyID/tags")
	{
		// POST endpoint to create a tag
		tag.POST("", baseInstance.CreateTag)

		// GET endpoint to retrieve all tag
		tag.GET("", baseInstance.ReadTagslist)

		// GET endpoint to retrieve details of a specific tag
		tag.GET("/:ID", baseInstance.ReadTag)

		// PUT endpoint to update details of a specific tag
		tag.PUT("/:ID", baseInstance.UpdateTag)

		// DELETE endpoint to delete a specific tag
		tag.DELETE("/:ID", baseInstance.DeleteTag)

		// Assign tags to mailinglist
		tag.POST("/:ID/mailinglist/:mailinglistID", baseInstance.AssignTagToMailinglist)

		//Assign tags to contact
		tag.POST("/:ID/contacts/:contactID", baseInstance.AssignTagToContact)

		//Delete tags from mailinglist
		tag.DELETE("/:ID/mailinglist/:mailinglistID", baseInstance.RemoveTagFromMailinglist)

		//Delete tags from contact
		tag.DELETE("/:ID/contacts/:contactID", baseInstance.RemoveTagFromContact)
	}
}
