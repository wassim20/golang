package campaign

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CampaignRouterInit initializes the routes related to companies.
func CampaignRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewCampaignRepository(db)

	// Private
	campaigns := router.Group("/:companyID/campaigns")
	{

		// POST endpoint to create a new campaign
		campaigns.POST("/:mailinglistID", baseInstance.CreateCampaign)

		// GET endpoint to retrieve all campaigns
		campaigns.GET("", baseInstance.ReadAllCampaigns)

		// GET endpoint to retrieve details of a specific campaign
		campaigns.GET("/:ID", baseInstance.ReadCampaign)

		// PUT endpoint to update details of a specific campaign
		campaigns.PUT("/:ID", baseInstance.UpdateCampaign)

		// DELETE endpoint to delete a specific campaign
		campaigns.DELETE("/:ID", baseInstance.DeleteCampaign)
	}

}
