package mailinglists

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CompanyRouterInit initializes the routes related to mailinglist.
func MailinglistRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewMailinglistRepository(db)

	// Private
	mailinglist := router.Group("/mailinglist")
	{

		// GET endpoint to retrieve all mailinglist
		mailinglist.GET("", baseInstance.ReadMailinglists)

		// GET endpoint to retrieve details of a specific company
		mailinglist.GET("/:ID", baseInstance.ReadMailinglist)

		// PUT endpoint to update details of a specific company
		//mailinglist.PUT("/:ID", baseInstance.UpdateCompany)

		// DELETE endpoint to delete a specific company
		//mailinglist.DELETE("/:ID", baseInstance.DeleteCompany)
	}
}
