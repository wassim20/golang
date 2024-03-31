package language

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LanguageRouterInit initializes the routes related to languages.
func LanguageRouterInit(router *gin.RouterGroup, db *gorm.DB) {

	// Initialize database instance
	baseInstance := Database{DB: db}

	// Automigrate / Update table
	NewLanguageRepository(db)

	// Private
	languages := router.Group("/languages")
	{

		// POST endpoint to create a new language
		languages.POST("", baseInstance.CreateLanguage)

		// GET endpoint to retrieve all languages
		languages.GET("", baseInstance.GetAllLanguages)

		// GET endpoint to retrieve details of a specific language
		languages.GET("/:ID", baseInstance.GetLanguageByID)

		// PUT endpoint to update details of a specific language
		languages.PUT("/:ID", baseInstance.UpdateLanguage)

		// DELETE endpoint to delete a specific language
		languages.DELETE("/:ID", baseInstance.DeleteLanguage)
	}
}
