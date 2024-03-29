package country

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the country package.
type Database struct {
	DB *gorm.DB
}

// NewCountryRepository performs automatic migration of country-related structures in the database.
func NewCountryRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Country{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the campaign structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of countries based on limit and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Country, limit, offset int) ([]domains.Country, error) {
	err := db.Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a country by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Country, id uuid.UUID) (domains.Country, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadTotalCount retrieves the total count of countries in the database.
func ReadTotalCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&domains.Country{}).Count(&count).Error
	return count, err
}
