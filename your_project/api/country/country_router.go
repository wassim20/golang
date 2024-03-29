package country

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CountryRouterInit initializes the routes related to countries.
func CountryRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewCountryRepository(db)

	// Private
	countries := router.Group("/countries")
	{

		// POST endpoint to create a new country
		countries.POST("", baseInstance.CreateCountry)

		// GET endpoint to retrieve all countries
		countries.GET("", baseInstance.GetAllCountries)

		// GET endpoint to retrieve details of a specific country
		countries.GET("/:ID", baseInstance.GetCountryByID)

		// PUT endpoint to update details of a specific country
		countries.PUT("/:ID", baseInstance.UpdateCountry)

		// DELETE endpoint to delete a specific country
		countries.DELETE("/:ID", baseInstance.DeleteCountry)
	}
}
