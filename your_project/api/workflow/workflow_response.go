package workflow

import (
	"encoding/json"
	"labs/domains"
	"time"

	"github.com/google/uuid"
)

// @Description	WorkflowsPagination represents the paginated list of workflows.
type WorkflowsPagination struct {
	Items      []WorkflowTable `json:"items"`      // Items is a slice containing individual workflow details.
	Page       uint            `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint            `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint            `json:"totalCount"` // TotalCount is the total number of workflows in the entire list.
} //@name WorkflowsPagination

// @Description	WorkflowIn represents the input structure for creating a new workflow.
type WorkflowIn struct {
	Name          string          `json:"name"`           // Name is the name of the workflow.
	Trigger       string          `json:"trigger"`        // Trigger is the trigger of the workflow.
	MailinglistID uuid.UUID       `json:"mailinglist_id"` // MailinglistID is the ID of the mailinglist associated with the workflow.
	Trigger_data  json.RawMessage `json:"trigger_data"`   // TriggerData is the data associated with the trigger.
} //@name WorkflowIn

// @Description	WorkflowDetails represents detailed information about a specific workflow.
type WorkflowDetails struct {
	ID            uuid.UUID        `json:"id"`             // ID is the unique identifier for the workflow.
	Name          string           `json:"name"`           // Name is the name of the workflow.
	CompanyID     uuid.UUID        `json:"company_id"`     // CompanyID is the ID of the company to which the workflow belongs.
	Status        string           `json:"status"`         // Status is the status of the workflow.
	CurrentStep   string           `json:"current_step"`   // CurrentStep is the current step of the workflow.
	Trigger       string           `json:"trigger"`        // Trigger is the trigger of the workflow.
	MailinglistID uuid.UUID        `json:"mailinglist_id"` // MailinglistID is the ID of the mailinglist associated with the workflow.
	Actions       []domains.Action `json:"actions"`        // Actions is a list of actions associated with the workflow.
	CreatedAt     time.Time        `json:"created_at"`     // CreatedAt is the timestamp indicating when the workflow entry was created.
} //@name WorkflowDetails

// @Description	WorkflowTable represents a single workflow entry in a table.
type WorkflowTable struct {
	ID            uuid.UUID `json:"id"`             // ID is the unique identifier for the workflow.
	Name          string    `json:"name"`           // Name is the name of the workflow.
	CompanyID     uuid.UUID `json:"company_id"`     // CompanyID is the ID of the company to which the workflow belongs.
	Status        string    `json:"status"`         // Status is the status of the workflow.
	CurrentStep   string    `json:"current_step"`   // CurrentStep is the current step of the workflow.
	Trigger       string    `json:"trigger"`        // Trigger is the trigger of the workflow.
	MailinglistID uuid.UUID `json:"mailinglist_id"` // MailinglistID is the ID of the mailinglist associated with the workflow.
	CreatedAt     time.Time `json:"created_at"`     // CreatedAt is the timestamp indicating when the workflow entry was created.
} //@name WorkflowTable
