package role

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	roleGroup := router.Group("/api/role")
	{
		roleGroup.POST("", CreateRole)
		roleGroup.GET("/:id", ReadRole)
		roleGroup.GET("", ReadAllRoles)
		roleGroup.PUT("/:id", UpdateRole)
		roleGroup.DELETE("/:id", DeleteRole)
		roleGroup.POST("/:role_id/user/:user_id", AssignUserToRole)
	}
}
