package role

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wassim_p/App/initializers"
	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

func CreateRole(c *gin.Context) {
	var service = NewRoleService(initializers.DB)
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Create(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create role: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func ReadRole(c *gin.Context) {
	var service = NewRoleService(initializers.DB)

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read role: %v", err)})
		return
	}

	role, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read role: %v", err)})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

func ReadAllRoles(c *gin.Context) {
	var service = NewRoleService(initializers.DB)

	roles, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read roles: %v", err)})
		return
	}

	c.JSON(http.StatusOK, roles)
}

func UpdateRole(c *gin.Context) {
	var service = NewRoleService(initializers.DB)
	var role models.Role

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read role: %v", err)})
		return
	}
	role.ID = uint(idUint)

	if err := c.ShouldBindJSON(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Update(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to update role: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, role)
}

func DeleteRole(c *gin.Context) {
	var service = NewRoleService(initializers.DB)
	var role models.Role

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete role: %v", err)})
		return
	}
	role.ID = uint(idUint)

	if err := service.Delete(&role); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete role: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, "Role deleted successfully")
}

func AssignUserToRole(c *gin.Context) {
	var service = NewRoleService(initializers.DB)

	roleID, err := strconv.ParseUint(c.Param("role_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse role ID: %v", err)})
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse user ID: %v", err)})
		return
	}

	err = service.AssignUser(uint(roleID), uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign user to role: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User assigned to role successfully"})
}
