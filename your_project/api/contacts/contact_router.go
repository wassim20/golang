package contacts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContactRouterInit initializes the routes related to companies.
func ContactRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewContactRepository(db)

	// Private
	contacts := router.Group("/contacts")
	{

		// GET endpoint to retrieve all contacts
		contacts.GET("", baseInstance.ReadContacts)

		// Post endpoint
		contacts.POST("", baseInstance.CreateContact)

		// GET endpoint to retrieve details of a specific company
		contacts.GET("/:ID", baseInstance.ReadContact)

		// PUT endpoint to update details of a specific company
		contacts.PUT("/:ID", baseInstance.UpdateContact)

		// DELETE endpoint to delete a specific company
		contacts.DELETE("/:ID", baseInstance.DeleteContact)
	}
}
