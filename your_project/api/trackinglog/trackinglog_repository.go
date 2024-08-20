package trackinglog

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the logs package.
type Database struct {
	DB *gorm.DB
}

// NewLogRepository performs automatic migration of log-related structures in the database.
func NewLogRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.TrackingLog{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the log structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of logs based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.TrackingLog, companyID uuid.UUID, campaignID uuid.UUID, limit, offset int) ([]domains.TrackingLog, error) {
	err := db.Where("company_id = ? and campaign_id = ?", companyID, campaignID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

func ReadAll(db *gorm.DB, model []domains.TrackingLog, companyID uuid.UUID, limit, offset int) ([]domains.TrackingLog, error) {
	err := db.Where("company_id = ? ", companyID).Find(&model).Error
	return model, err
}

// ReadByID retrieves a log by its unique identifier.
func ReadByID(db *gorm.DB, model domains.TrackingLog, companyID uuid.UUID, campaignID uuid.UUID, id uuid.UUID) (domains.TrackingLog, error) {
	err := db.Where("company_id = ? and campaign_id = ?", companyID, campaignID).First(&model, id).Error
	return model, err
}

func ReadTotalCountTrackingLog(db *gorm.DB, companyID uuid.UUID, campaignID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&domains.TrackingLog{}).Where("company_id = ? and campaign_id = ?", companyID, campaignID).Count(&count).Error
	return count, err
}
func ReadTotalCountAllTrackingLog(db *gorm.DB, companyID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&domains.TrackingLog{}).Where("company_id = ? ", companyID).Count(&count).Error
	return count, err
}
