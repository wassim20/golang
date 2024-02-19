package tags

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

// NewTagRepository performs automatic migration of tag-related structures in the database.
func NewTagRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Tag{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the tag structure. Error: ", err)
	}
}

func ReadAllPagination(db *gorm.DB, model []domains.Tag, modelID uuid.UUID, limit, offset int) ([]domains.Tag, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}
