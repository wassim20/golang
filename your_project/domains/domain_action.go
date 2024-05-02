package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Action struct {
	ID         uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the action
	Name       string    `gorm:"column:name; not null;"`                      // The action's name
	ParentID   uuid.UUID `gorm:"column:parent_id; type:uuid; not null;"`      // ID of the parent action
	WorkflowID uuid.UUID `gorm:"column:workflow_id; type:uuid; not null;"`    // ID of the workflow to which the action belongs
	Type       string    `gorm:"column:type; not null;"`                      // The action's type
	Status     string    `gorm:"column:status; not null;"`                    // The action's status
	// add a field to sttore the data needed to perform the action wich is a json field
	Data string `gorm:"column:data; not null;"` // The action's data
	gorm.Model
}

func ReadActionNameByID(db *gorm.DB, actionID uuid.UUID) (string, error) {
	action := new(Action)
	check := db.Select("name").Where("id = ?", actionID).First(action)

	if check.Error != nil {
		return "", check.Error
	}

	return action.Name, nil
}
