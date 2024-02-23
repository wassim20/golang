package api

import (
	"labs/api/auth"
	"labs/api/companies"
	"labs/api/mailinglists"
	"labs/api/notifications"
	"labs/api/roles"
	"labs/api/tags"
	"labs/api/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoutesApiInit initializes the API routes for various modules.
func RoutesApiInit(router *gin.Engine, db *gorm.DB) {

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
	}
}
