package company

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	companyGroup := router.Group("/api/company")
	{
		// CRUD routes
		companyGroup.POST("", CreateCompany)
		companyGroup.GET("/:id", ReadCompany)
		companyGroup.GET("", ReadAllCompanies)
		companyGroup.PUT("/:id", UpdateCompany)
		companyGroup.DELETE("/:id", DeleteCompany)

		// Placeholder for future routes (e.g., assigning mailing lists/tags)
		companyGroup.POST("/:id/mailinglist/:mailinglist_id", AssignMailingListToCompany)
		companyGroup.POST("/:id/tag/:tag_id", AssignTagToCompany)
		companyGroup.POST("/:id/user/:user_id", AssignUserToCompany)
	}
}
