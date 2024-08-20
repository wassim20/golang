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
	dashboardTrackinglog := router.Group("/:companyID/logs")
	{
		// GET endpoint to get trackinglogs for updateChartData
		dashboardTrackinglog.GET("barchartdata", baseInstance.updateChartData)
		// GET endpoint to get trackinglogs for updatePieChartData
		dashboardTrackinglog.GET("piechartdata", baseInstance.updatePieChartData)
		// GET endpoint to get trackinglogs for updateRadialChartData
		dashboardTrackinglog.GET("radialchartdata", baseInstance.updateRadialChartData)
		// GET endpoint to get trackinglogs for updateLineChartData
		dashboardTrackinglog.GET("linechartdata", baseInstance.updateLineChartData)
		// GET endpoint to get trackinglogs for updateScatterChartData
		dashboardTrackinglog.GET("scatterchartdata", baseInstance.updateScatterChartData)
		// GET endpoint to get trackinglogs for barChartOpens
		dashboardTrackinglog.GET("barchartopens", baseInstance.barChartDataOpens)
		// GET endpoint to get trackinglogs for barChartClicks
		dashboardTrackinglog.GET("barchartclicks", baseInstance.barChartDataClicks)
		// GET endpoint to get trackinglogs for ScatterChartOpens
		dashboardTrackinglog.GET("scatterchartopens", baseInstance.scatterChartDataOpens)
		//GET endpoint to get trackinglogs for ScatterChartClicks
		dashboardTrackinglog.GET("scatterchartclicks", baseInstance.scatterChartDataClicks)

	}

	trackinglogs := router.Group("/:companyID/:campaignID/logs")
	{

		// POST endpoint to create a new TrackingLog
		trackinglogs.POST("", baseInstance.CreateTrackingLog)

		// GET endpoint to retrieve all trackinglogs
		trackinglogs.GET("", baseInstance.ReadTrackingLogs)

		// GET endpoint to retrieve details of a specific tracking
		trackinglogs.GET("/:ID", baseInstance.ReadTrackingLogByID)

		// PUT endpoint to update details of a specific tracking
		trackinglogs.PUT("/:ID", baseInstance.UpdateTrackingLog)

		// DELETE endpoint to delete a specific tracking
		trackinglogs.DELETE("/:ID", baseInstance.DeleteTrackingLog)

		//POST endpoint to update trackinglog when link is clicked
		trackinglogs.POST("/click/:trackingID/:email", baseInstance.handleClickRequest)

	}

	trackinglogsworkflow := router.Group("/:companyID/logs")
	{
		// GET endpoint to retrieve all trackinglogs for all campaigns
		trackinglogsworkflow.GET("", baseInstance.ReadAllTrackingLogs)
		// POST endpoint to handle the opening of an email
		trackinglogsworkflow.POST("/open/:trackingID", baseInstance.handleOpenRequestWorflow)

		// POST endpoint to handle the clicking of a link
		trackinglogsworkflow.POST("/click/:trackingID", baseInstance.handleClickRequestWorflow)

	}
	trackingpixel := router.Group("/track")
	{
		trackingpixel.GET("/pixel.png", baseInstance.handleOpenRequest)
	}

}
