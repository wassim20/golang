package action

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActionRouterInit initializes the routes related to action.
func ActionRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewActionRepository(db)

	// Private
	action := router.Group("/:workflowID/action")
	{
		//POST endpoint to create an action
		action.POST("", baseInstance.CreateAction)

		// GET endpoint to retrieve all actions
		action.GET("", baseInstance.ReadActions)

		// GET endpoint to retrieve details of a specific action
		action.GET("/:actionID", baseInstance.ReadAction)

		// PUT endpoint to update details of a specific action
		action.PUT("/:actionID", baseInstance.UpdateAction)

		// DELETE endpoint to delete a specific action
		action.DELETE("/:actionID", baseInstance.DeleteAction)

	}

}
