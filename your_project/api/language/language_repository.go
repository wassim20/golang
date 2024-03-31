package language

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

// NewLanguageRepository performs automatic migration of language-related structures in the database.
func NewLanguageRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Language{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the language structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of languages based on limit and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Language, limit, offset int) ([]domains.Language, error) {
	err := db.Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a language by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Language, id uuid.UUID) (domains.Language, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadTotalCount retrieves the total count of languages in the database.
func ReadTotalCount(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&domains.Language{}).Count(&count).Error
	return count, err
}
