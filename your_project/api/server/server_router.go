package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServerRouterInit initializes the routes related to servers.
func ServerRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewServerRepository(db)

	// Private
	servers := router.Group("/:companyID/servers")
	{

		// GET endpoint to retrieve all servers
		servers.GET("", baseInstance.ReadServers)

		// Post endpoint
		servers.POST("", baseInstance.CreateServer)

		// GET endpoint to retrieve details of a specific server
		servers.GET("/:ID", baseInstance.ReadServer)

		// PUT endpoint to update details of a specific server
		servers.PUT("/:ID", baseInstance.UpdateServer)

		// DELETE endpoint to delete a specific server
		servers.DELETE("/:ID", baseInstance.DeleteServer)

	}

}
