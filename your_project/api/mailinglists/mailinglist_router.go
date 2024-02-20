package mailinglists

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MailinglistRouterInit initializes the routes related to mailinglist.
func MailinglistRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewMailinglistRepository(db)

	// Private
	mailinglist := router.Group("/mailinglist")
	{
		//POST endpoint ro create a mailinglist
		mailinglist.POST("", baseInstance.CreateMailinglist)

		// GET endpoint to retrieve all mailinglist
		mailinglist.GET("", baseInstance.ReadMailinglists)

		// GET endpoint to retrieve details of a specific mailinglist
		mailinglist.GET("/:ID", baseInstance.ReadMailinglist)

		// PUT endpoint to update details of a specific mailinglist
		mailinglist.PUT("/:ID", baseInstance.UpdateMailinglist)

		// DELETE endpoint to delete a specific mailinglist
		mailinglist.DELETE("/:ID", baseInstance.DeleteMailinglist)
	}
}
