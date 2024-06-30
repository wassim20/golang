package api

import (
	"labs/api/auth"
	"labs/api/campaign"
	"labs/api/companies"
	"labs/api/country"
	"labs/api/language"
	"labs/api/mailinglists"
	"labs/api/notifications"
	"labs/api/roles"
	"labs/api/server"
	"labs/api/tags"
	"labs/api/trackinglog"
	"labs/api/users"
	"labs/api/workflow"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoutesApiInit initializes the API routes for various modules.
func RoutesApiInit(router *gin.Engine, db *gorm.DB) {

	// CORS middleware configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Replace with your Angular app's URL
	config.AllowHeaders = []string{"Authorization"}
	router.Use(cors.New(config))
	api := router.Group("/api")
	{
		// Initialize authentication routes
		auth.AuthRouterInit(api, db)

		// Initialize user routes
		users.UserRouterInit(api, db)

		// Initialize company routes
		companies.CompanyRouterInit(api, db)

		// Initialize role routes
		roles.RoleRouterInit(api, db)

		// Initialize notification routes
		notifications.NotificationRouterInit(api, db)

		// Initialize mailinglist routes
		mailinglists.MailinglistRouterInit(api, db)

		// Initialize contact routes
		//contacts.ContactRouterInit(api, db)

		// Initialize tag routes
		tags.TagRouterInit(api, db)

		// Initialize campaign routes
		campaign.CampaignRouterInit(api, db)

		// Initialize trackinglog routes
		trackinglog.TrackingLogRouterInit(api, db)

		// Initialize language routes
		language.LanguageRouterInit(api, db)

		// Inistialize country routes
		country.CountryRouterInit(api, db)

		// Initialize email server routes
		server.ServerRouterInit(api, db)

		// Initialize workflow routes
		workflow.WorkflowRouterInit(api, db)

		api.Static("static", "./static")

	}
}
