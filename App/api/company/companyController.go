package company

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

type Companyout struct {
	Email              string
	Password           string
	AccountType        string
	SubscriptionOption string

	MailingLists []models.MailingList
	Tags         []models.Tag
}

func CreateCompany(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)
	var company models.Company

	// 1. Parse request body and validate data
	if err := c.ShouldBindJSON(&company); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 2. Call companyService.Create to create the company in the database

	err := service.Create(&company)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create company: %v", err)})
		return
	}
	fmt.Println("hellloooo")
	// 3. Return successful response with the created company data
	c.JSON(http.StatusCreated, company)
}

func ReadCompany(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)

	ID := c.Param("id")

	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read company: %v", err)})
		return
	}
	fmt.Println(idUint)

	company, err := service.Read(uint(idUint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read company: %v", err)})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)

}

func ReadAllCompanies(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)

	companies, err := service.ReadAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read companies: %w", err)})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func UpdateCompany(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)
	var company models.Company
	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to read company: %v", err)})
		return
	}
	company.ID = uint(idUint)

	if err := c.ShouldBindJSON(&company); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if service.Update(&company); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to create company: %v", err)})
		return
	}
	fmt.Println("hellloooo")
	// 3. Return successful response with the created company data
	c.JSON(http.StatusCreated, company)

}

func DeleteCompany(c *gin.Context) {

	var service = NewCompanyService(initializers.DB)
	var company models.Company
	ID := c.Param("id")
	idUint, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete company: %v", err)})
		return
	}
	company.ID = uint(idUint)

	if service.Delete(&company); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to delete company: %v", err)})
		return
	}

	// 3. Return successful response with the created company data
	c.JSON(http.StatusAccepted, "Company Deleted successfully")
}

func AssignMailingListToCompany(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)

	companyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse company ID: %v", err)})
		return
	}

	mailingListID, err := strconv.ParseUint(c.Param("mailinglist_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse mailing list ID: %v", err)})
		return
	}

	err = service.AssignMailingList(uint(companyID), uint(mailingListID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign mailing list to company: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mailing list assigned to company successfully"})
}

func AssignTagToCompany(c *gin.Context) {
	var service = NewCompanyService(initializers.DB)

	companyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse company ID: %v", err)})
		return
	}

	tagID, err := strconv.ParseUint(c.Param("tag_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse tag ID: %v", err)})
		return
	}

	err = service.AssignTag(uint(companyID), uint(tagID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign tag to company: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag assigned to company successfully"})
}

func AssignUserToCompany(c *gin.Context) {

	var service = NewCompanyService(initializers.DB)

	companyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse company ID: %v", err)})
		return
	}

	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to parse tag ID: %v", err)})
		return
	}

	err = service.AssignUser(uint(companyID), uint(userID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("Failed to assign user to company: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user assigned to company successfully"})

}
