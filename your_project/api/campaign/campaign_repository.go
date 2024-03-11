package campaign

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the companies package.
type Database struct {
	DB *gorm.DB
}

// NewCampaignRepository performs automatic migration of campaign-related structures in the database.
func NewCampaignRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Campaign{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the campaign structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of campaigns based on mailinglist_ID, limit, and offset.
func ReadAllPaginationFromMailinglist(db *gorm.DB, model []domains.Campaign, modelID uuid.UUID, limit, offset int) ([]domains.Campaign, error) {
	err := db.Where("mailinglist_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a campaign by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Campaign, id uuid.UUID) (domains.Campaign, error) {
	err := db.First(&model, id).Error
	return model, err
}
