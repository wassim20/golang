package trackinglog

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TrackingLogRouterInit initializes the routes related to TrackingLog.
func TrackingLogRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewLogRepository(db)

	// Private

	trackinglogs := router.Group("/:companyID/:camapignID/logs/")
	{

		// POST endpoint to create a new TrackingLog
		trackinglogs.POST("", baseInstance.CreateTrackingLog)

		// GET endpoint to retrieve all trackinglogs
		trackinglogs.GET("", baseInstance.ReadTrackingLogs)

		// GET endpoint to retrieve details of a specific company
		trackinglogs.GET("/:ID", baseInstance.ReadTrackingLogByID)

		// PUT endpoint to update details of a specific company
		trackinglogs.PUT("/:ID", baseInstance.UpdateTrackingLog)

		// DELETE endpoint to delete a specific company
		trackinglogs.DELETE("/:ID", baseInstance.DeleteTrackingLog)

		//POST endpoint to update trackinglog when email is opened
		trackinglogs.POST("/open/:trackingID", baseInstance.handleOpenRequest)

		//POST endpoint to update trackinglog when link is clicked
		trackinglogs.POST("/click/:trackingID/:email", baseInstance.handleClickRequest)
	}

	trackinglogsworkflow := router.Group("/:companyID/logs")
	{
		// POST endpoint to handle the opening of an email
		trackinglogsworkflow.POST("/open/:trackingID", baseInstance.handleOpenRequestWorflow)

		// POST endpoint to handle the clicking of a link
		trackinglogsworkflow.POST("/click/:trackingID", baseInstance.handleClickRequestWorflow)

	}
}
