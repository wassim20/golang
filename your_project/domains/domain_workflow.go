/*
Package domains provides the data structures representing entities in the project.

Structures:
- Workflow: Represents information about a workflow in the system.
  - ID (uuid.UUID): Unique identifier for the workflow.
  - Name (string): The name of the workflow.
  - Status (string): The workflow's status.
  - CurrentStep (string): The workflow's current step.
  - Trigger (string): The workflow's trigger.
  - MailinglistID (uuid.UUID): ID of the mailinglist associated with the workflow.
  - Actions ([]Action): List of actions associated with the workflow.
  - gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

Functions:
- ReadWorkflowNameByID(db *gorm.DB, workflowID uuid.UUID) (string, error): Retrieves the name of the workflow based on its ID from the database.

Dependencies:
- "github.com/google/uuid": Package for working with UUIDs.
- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

Usage:
- Import this package to utilize the provided data structures and functions for handling entities in the project.

Note:
- The Workflow structure represents information about a workflow in the system.
- ReadWorkflowNameByID retrieves the name of the workflow based on its ID from the database.

Last update :
29/04/2024 6:10
*/
package domains

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Workflow struct {
	ID            uuid.UUID       `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the workflow
	Name          string          `gorm:"column:name; not null;"`                      // The workflow's name
	CompanyID     uuid.UUID       `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the worflow belongs
	Status        string          `gorm:"column:status; not null;"`                    // The workflow's status
	CurrentStep   string          `gorm:"column:current_step; not null;"`              // The workflow's current step
	Trigger       string          `gorm:"column:trigger; not null;"`
	TriggerData   json.RawMessage `gorm:"column:trigger_data; type:json; default:'{}'"`                                        // The workflow's trigger
	MailinglistID uuid.UUID       `gorm:"column:mailinglist_id; type:uuid; not null;"`                                         // ID of the mailinglist associated with the workflow
	Actions       []Action        `gorm:"foreignKey:workflow_id; references:id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // List of actions associated with the workflow
	gorm.Model
}

func ReadWorkflowNameByID(db *gorm.DB, workflowID uuid.UUID) (string, error) {
	workflow := new(Workflow)
	check := db.Select("name").Where("id = ?", workflowID).First(workflow)

	if check.Error != nil {
		return "", check.Error
	}

	return workflow.Name, nil
}
