package action

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

// NewActionRepository performs automatic migration of action-related structures in the database.
func NewActionRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Action{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the action structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of actions based on workflow ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Action, modelID uuid.UUID, limit, offset int) ([]domains.Action, error) {
	err := db.Where("workflow_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves an action by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Action, workflowID uuid.UUID, id uuid.UUID) (domains.Action, error) {
	err := db.Where("workflow_id = ?", workflowID).First(&model, id).Error
	return model, err
}

func ReadAllList(db *gorm.DB, model []domains.Action, modelID uuid.UUID) ([]domains.Action, error) {
	err := db.Where("workflow_id = ? ", modelID).Find(&model).Error
	return model, err
}
