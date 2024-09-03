package action

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// @Description	ActionsPagination represents the paginated list of actions.
type ActionsPagination struct {
	Items      []ActionTable `json:"items"`      // Items is a slice containing individual action details.
	Page       uint          `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint          `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint          `json:"totalCount"` // TotalCount is the total number of actions in the entire list.
} //@name ActionsPagination

// @Description	ActionIn represents the input structure for creating a new action.
type ActionIn struct {
	Name     string          `json:"name"`      // Name is the name of the action.
	Type     string          `json:"type"`      // Type is the type of the action.
	Data     json.RawMessage `json:"data"`      // Data is the data needed to perform the action.
	ParentID uuid.UUID       `json:"parent_id"` // ParentID is the ID of the parent action.
} //@name ActionIn
type Actionup struct {
	Type string          `json:"type"` // Type is the type of the action.
	Data json.RawMessage `json:"data"` // Data is the data needed to perform the action.
}

// @Description	ActionDetails represents detailed information about a specific action.
type ActionDetails struct {
	ID         uuid.UUID       `json:"id"`          // ID is the unique identifier for the action.
	Name       string          `json:"name"`        // Name is the name of the action.
	ParentID   uuid.UUID       `json:"parent_id"`   // ParentID is the ID of the parent action.
	Type       string          `json:"type"`        // Type is the type of the action.
	Status     string          `json:"status"`      // Status is the status of the action.
	WorkflowID uuid.UUID       `json:"workflow_id"` // WorkflowID is the ID of the workflow to which the action belongs.
	Data       json.RawMessage `json:"data"`        // Data is the data needed to perform the action.
	CreatedAt  time.Time       `json:"created_at"`  // CreatedAt is the timestamp indicating when the action entry was created.
} //@name ActionDetails

// @Description	ActionTable represents a single action entry in a table.
type ActionTable struct {
	ID         uuid.UUID       `json:"id"`          // ID is the unique identifier for the action.
	Name       string          `json:"name"`        // Name is the name of the action.
	Type       string          `json:"type"`        // Type is the type of the action.
	ParentID   uuid.UUID       `json:"parent_id"`   // ParentID is the ID of the parent action.
	Status     string          `json:"status"`      // Status is the status of the action.
	WorkflowID uuid.UUID       `json:"workflow_id"` // WorkflowID is the ID of the workflow to which the action belongs.
	Data       json.RawMessage `json:"data"`        // Data is the data needed to perform the action.
	CreatedAt  time.Time       `json:"created_at"`  // CreatedAt is the timestamp indicating when the action entry was created.
} //@name ActionTable
