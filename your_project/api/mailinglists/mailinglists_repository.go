package mailinglists

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the mailinglist package.
type Database struct {
	DB *gorm.DB
}

// NewmailinglistRepository performs automatic migration of mailinglist-related structures in the database.
func NewMailinglistRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Mailinglist{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the mailinglist structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of mailinglist based on mailinglist ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Mailinglist, modelID uuid.UUID, limit, offset int) ([]domains.Mailinglist, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a mailinglist by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Mailinglist, id uuid.UUID) (domains.Mailinglist, error) {
	err := db.First(&model, id).Error
	return model, err
}

func ReadAllList(db *gorm.DB, model []domains.Mailinglist, modelID uuid.UUID) ([]domains.Mailinglist, error) {
	err := db.Where("company_id = ? ", modelID).Find(&model).Error
	return model, err
}
