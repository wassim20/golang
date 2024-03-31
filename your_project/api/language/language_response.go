package language

import (
	"time"

	"github.com/google/uuid"
)

// @Description	 LanguagePagination represents the paginated list of languages.
type LanguagePagination struct {
	Items      []LanguageTable `json:"items"`      // Items is a slice containing individual language details.
	Page       uint            `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint            `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint            `json:"totalCount"` // TotalCount is the total number of languages in the entire list.
} //@name LanguagePagination

// @Description	 LanguageIn represents the input structure for creating a new language.
type LanguageIn struct {
	Name string `json:"name"` // Name is the name of the language.
	Code string `json:"code"` // Code is the code of the language.
} //@name LanguageIn

// @Description	 LanguageDetails represents detailed information about a specific language.
type LanguageDetails struct {
	ID        uuid.UUID  `json:"id"`         // ID is the unique identifier for the language.
	Name      string     `json:"name"`       // Name is the name of the language.
	Code      string     `json:"code"`       // Code is the code of the language.
	CreatedAt time.Time  `json:"created_at"` // CreatedAt is the timestamp indicating when the language entry was created.
	UpdatedAt time.Time  `json:"updated_at"` // UpdatedAt is the timestamp indicating when the language entry was updated.
	DeletedAt *time.Time `json:"deleted_at"` // DeletedAt is the timestamp indicating when the language entry was deleted.
} //@name LanguageDetails

// @Description	 LanguageTable represents a single language entry in a table.
type LanguageTable struct {
	Name string `json:"name"` // Name is the name of the language.
	Code string `json:"code"` // Code is the code of the language.
} //@name LanguageTable
