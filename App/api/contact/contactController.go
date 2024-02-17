package contact

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

func CreateContact(c *gin.Context) {
	var service = NewContactService(initializers.DB)
	var contact models.Contact

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Create(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create contact: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func ReadContact(c *gin.Context) {
	var service = NewContactService(initializers.DB)

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read contact: %v", err)})
		return
	}

	contact, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read contact: %v", err)})
		return
	}

	if contact == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	c.JSON(http.StatusOK, contact)
}

func ReadAllContacts(c *gin.Context) {
	var service = NewContactService(initializers.DB)

	contacts, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read contacts: %v", err)})
		return
	}

	c.JSON(http.StatusOK, contacts)
}

func UpdateContact(c *gin.Context) {
	var service = NewContactService(initializers.DB)
	var contact models.Contact

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read contact: %v", err)})
		return
	}
	contact.ID = uint(idUint)

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := service.Update(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to update contact: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func DeleteContact(c *gin.Context) {
	var service = NewContactService(initializers.DB)
	var contact models.Contact

	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete contact: %v", err)})
		return
	}
	contact.ID = uint(idUint)

	if err := service.Delete(&contact); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete contact: %v", err)})
		return
	}

	c.JSON(http.StatusAccepted, "Contact deleted successfully")
}

func AssignTagToContact(c *gin.Context) {
	var service = NewContactService(initializers.DB)

	contactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse contact ID: %v", err)})
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse tag ID: %v", err)})
		return
	}

	err = service.AssignTag(uint(contactID), uint(tagID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign tag to contact: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag assigned to contact successfully"})
}
