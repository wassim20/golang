package tags

import (
	"time"

	"github.com/google/uuid"
)

// @Description TagPaginator represents the paginated list of tags.
type TagPaginator struct {
	Items      []TagTable `json:"items"`      // Items is a slice containing individual tag details.
	Page       uint       `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint       `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint       `json:"totalCount"` // TotalCount is the total number of tags in the entire list.
}

// @Description TagIn represents the input structure for creating a new tag.
type TagIn struct {
	Name string `json:"name"` // Tag's name.
}

// @Description TagDetails represents detailed information about a specific tag.
type TagDetails struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the tag.
	Name      string    `json:"name"`       // Tag's  name.
	CompanyID uuid.UUID `json:"company_id"` // Tag's company
	CreatedAt time.Time `json:"created_at"` // CreatedAt is the timestamp indicating when the tag was created.
}

// @Description TagTable represents a single tag entry in a table.
type TagTable struct {
	ID   uuid.UUID `json:"id"`   // Unique identifier for the tag.
	Name string    `json:"name"` // Tag's full name.
}
