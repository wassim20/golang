package tag

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

func CreateTag(c *gin.Context) {
	var service = NewTagService(initializers.DB)
	var tag models.Tag

	// Parse request body and validate data
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call tagService.Create to create the tag in the database
	if err := service.Create(&tag); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create tag: %v", err)})
		return
	}

	// Return successful response with the created tag data
	c.JSON(http.StatusCreated, tag)
}

func ReadTag(c *gin.Context) {
	var service = NewTagService(initializers.DB)

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read tag: %v", err)})
		return
	}

	tag, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read tag: %v", err)})
		return
	}

	if tag == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	c.JSON(http.StatusOK, tag)
}

func ReadAllTags(c *gin.Context) {
	var service = NewTagService(initializers.DB)

	tags, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read tags: %v", err)})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func UpdateTag(c *gin.Context) {
	var service = NewTagService(initializers.DB)
	var tag models.Tag

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read tag: %v", err)})
		return
	}
	tag.ID = uint(idUint)

	if err := c.ShouldBindJSON(&tag); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Update(&tag); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to update tag: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

func DeleteTag(c *gin.Context) {
	var service = NewTagService(initializers.DB)
	var tag models.Tag

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete tag: %v", err)})
		return
	}
	tag.ID = uint(idUint)

	if err := service.Delete(&tag); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete tag: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, "Tag deleted successfully")
}
