package mailinglist

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

func CreateMailinglist(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)
	var mailinglist models.MailingList

	if err := c.ShouldBindJSON(&mailinglist); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := service.Create(&mailinglist)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create mailinglist: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, mailinglist)
}

func ReadMailinglist(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)

	ID := c.Param("id")

	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read mailinglist: %v", err)})
		return
	}
	fmt.Println(idUint)

	mailinglist, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "mailinglist not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read company: %v", err)})
		return
	}

	if mailinglist == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, mailinglist)

}

func ReadAllMailinglists(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)

	mailinglists, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read mailinglists: %w", err)})
		return
	}

	c.JSON(http.StatusOK, mailinglists)
}

func UpdateMailinglist(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)
	var mailinglist models.MailingList
	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read mailinglist: %v", err)})
		return
	}
	mailinglist.ID = uint(idUint)

	if err := c.ShouldBindJSON(&mailinglist); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if service.Update(&mailinglist); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create mailinglist: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, mailinglist)

}

func DeleteMailinglist(c *gin.Context) {

	var service = NewMailingListService(initializers.DB)
	var mailinglist models.MailingList
	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete mailinglist: %v", err)})
		return
	}
	mailinglist.ID = uint(idUint)

	if service.Delete(&mailinglist); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete mailinglist: %v", err)})
		return
	}

	// 3. Return successful response with the created mailinglist data
	c.JSON(http.StatusAccepted, "mailinglist Deleted successfully")
}

func AssignContactToMailingList(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)

	mailingListID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse mailing list ID: %v", err)})
		return
	}

	contactID, err := strconv.ParseUint(c.Param("contact_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse contact ID: %v", err)})
		return
	}

	err = service.AssignContact(uint(mailingListID), uint(contactID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign contact to mailing list: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact assigned to mailing list successfully"})
}

func AssignTagToMailingList(c *gin.Context) {
	var service = NewMailingListService(initializers.DB)

	mailingListID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse mailing list ID: %v", err)})
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse tag ID: %v", err)})
		return
	}

	err = service.AssignTag(uint(mailingListID), uint(tagID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign tag to mailing list: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag assigned to mailing list successfully"})
}
