package workflow

import (
	"labs/api/action"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WorkflowRouterInit initializes the routes related to workflow.
func WorkflowRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update tablev
	NewWorkflowRepository(db)

	// Private
	workflow := router.Group("/:companyID/workflow")
	{
		//POST endpoint to create a workflow
		workflow.POST("", baseInstance.CreateWorkflow)

		// GET endpoint to retrieve all workflows
		workflow.GET("", baseInstance.ReadWorkflows)

		// GET endpoint to retrieve details of a specific workflow
		workflow.GET("/:workflowID", baseInstance.ReadWorkflow)

		// PUT endpoint to update details of a specific workflow
		workflow.PUT("/:workflowID", baseInstance.UpdateWorkflow)

		// DELETE endpoint to delete a specific workflow
		workflow.DELETE("/:workflowID", baseInstance.DeleteWorkflow)

		// POST endpoint to start a workflow
		workflow.POST("/:workflowID/start", baseInstance.StartWorkflow)

		action.ActionRouterInit(workflow, db)

	}

}
