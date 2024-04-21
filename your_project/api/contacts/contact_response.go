package contacts

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// @Description ContactPaginator represents the paginated list of contacts.
type ContactPaginator struct {
	Items      []ContactTable `json:"items"`      // Items is a slice containing individual contact details.
	Page       uint           `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint           `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint           `json:"totalCount"` // TotalCount is the total number of contacts in the entire list.
} //@name ContactPaginator

// @Description ContactIn represents the input structure for creating a new contact.
type ContactIn struct {
	Email       string `json:"email"`        // Contact's email address.
	Firstname   string `json:"first_name"`   // Contact's first name.
	Lastname    string `json:"last_name"`    // Contact's last name.
	PhoneNumber string `json:"phone_number"` // Contact's phone number.
	FullName    string `json:"full_name"`    // Contact's full name.
} //@name ContactIn

// @Description ContactDetails represents detailed information about a specific contact.
type ContactDetails struct {
	ID          uuid.UUID      `json:"id"`           // Unique identifier for the contact.
	Email       string         `json:"email"`        // Contact's email address.
	Firstname   string         `json:"first_name"`   // Contact's first name.
	Lastname    string         `json:"last_name"`    // Contact's last name.
	PhoneNumber string         `json:"phone_number"` // Contact's phone number.
	FullName    string         `json:"full_name"`    // Contact's full name.
	Tags        pq.StringArray `json:"tags"`         // Tags is a list of tag's UUID associated with the contact.
	CreatedAt   time.Time      `json:"created_at"`   // CreatedAt is the timestamp indicating when the contact was created.
} //@name ContactDetails

// @Description ContactTable represents a single contact entry in a table.
type ContactTable struct {
	ID          uuid.UUID `json:"id"`           // Unique identifier for the contact.
	Email       string    `json:"email"`        // Contact's email address.
	Firstname   string    `json:"first_name"`   // Contact's first name.
	Lastname    string    `json:"last_name"`    // Contact's last name.
	PhoneNumber string    `json:"phone_number"` // Contact's phone number.
	FullName    string    `json:"full_name"`    // Contact's full name.
	CreatedAt   time.Time `json:"created_at"`   // CreatedAt is the timestamp indicating when the contact was created.
} //@name ContactTable
