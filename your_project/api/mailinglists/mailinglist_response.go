package mailinglists

import (
	"labs/domains"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description	MailinglistsPagination represents the paginated list of mailinglists.
type MailinglistsPagination struct {
	Items      []MailinglistTable `json:"items"`      // Items is a slice containing individual mailinglist details.
	Page       uint               `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint               `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint               `json:"totalCount"` // TotalCount is the total number of mailinglists in the entire list.
} //@name MailinglistsPagination

// @Description	MailinglistIn represents the input structure for creating a new mailinglist.
type MailinglistIn struct {
	Name        string `json:"name"`        // Name is the name of the mailinglist.
	Description string `json:"description"` // Description is the description of the mailinglist.

} //@name MailinglistIn

// @Description	MailinglistDetails represents detailed information about a specific mailinglist.
type MailinglistDetails struct {
	ID              uuid.UUID         `json:"id"`              // ID is the unique identifier for the mailinglist.
	Name            string            `json:"name"`            // Name is the name of the mailinglist.
	Description     string            `json:"description"`     // Description is the description of the mailinglist.
	CompanyID       uuid.UUID         `json:"company_id"`      // CompanyID is the ID of the company to which the mailinglist belongs.
	Contacts        []domains.Contact `json:"contacts"`        // Contacts is a list of contacts associated with the mailinglist.
	Tags            pq.StringArray    `json:"tags"`            // Tags is a list of tag's UUID associated with the mailinglist.
	CreatedByUserID uuid.UUID         `json:"created_by_user"` // CreatedByUserID is the ID of the user who created the mailinglist.
	CreatedAt       time.Time         `json:"created_at"`      // CreatedAt is the timestamp indicating when the mailinglist entry was created.
} //@name MailinglistDetails

// @Description	MailinglistTable represents a single mailinglist entry in a table.
type MailinglistTable struct {
	ID              uuid.UUID `json:"id"`              // ID is the unique identifier for the mailinglist.
	Name            string    `json:"name"`            // Name is the name of the mailinglist.
	Description     string    `json:"description"`     // Description is the description of the mailinglist.
	CompanyID       uuid.UUID `json:"company_id"`      // CompanyID is the ID of the company to which the mailinglist belongs.
	CreatedByUserID uuid.UUID `json:"created_by_user"` // CreatedByUserID is the ID of the user who created the mailinglist.
	CreatedAt       time.Time `json:"created_at"`      // CreatedAt is the timestamp indicating when the mailinglist entry was created.
} //@name MailinglistTable
