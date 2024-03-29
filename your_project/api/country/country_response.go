package country

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Country represents the structure of a country.
type Country struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the country
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Code      string    `json:"code" gorm:"type:varchar(255);not null"`
	Currency  string    `json:"currency" gorm:"type:varchar(255);not null"`
	PhoneCode string    `json:"phone_code" gorm:"type:varchar(255);not null"`
	Flag      string    `json:"flag" gorm:"type:varchar(255);not null"`
	gorm.Model
}

// @Description	CountryPagination represents the paginated list of countries.
type CountryPagination struct {
	Items      []CountryTable `json:"items"`      // Items is a slice containing individual country details.
	Page       uint           `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint           `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint           `json:"totalCount"` // TotalCount is the total number of countries in the entire list.
} //@name CountryPagination

// @Description	CountryIn represents the input structure for creating a new country.
type CountryIn struct {
	Name      string `json:"name"`       // Name is the name of the country.
	Code      string `json:"code"`       // Code is the code of the country.
	Currency  string `json:"currency"`   // Currency is the currency of the country.
	PhoneCode string `json:"phone_code"` // PhoneCode is the phone code of the country.
	Flag      string `json:"flag"`       // Flag is the flag of the country.
} //@name CountryIn

// @Description	CountryDetails represents detailed information about a specific country.
type CountryDetails struct {
	ID        uuid.UUID `json:"id"`         // ID is the unique identifier for the country.
	Name      string    `json:"name"`       // Name is the name of the country.
	Code      string    `json:"code"`       // Code is the code of the country.
	Currency  string    `json:"currency"`   // Currency is the currency of the country.
	PhoneCode string    `json:"phone_code"` // PhoneCode is the phone code of the country.
	Flag      string    `json:"flag"`       // Flag is the flag of the country.
	CreatedAt time.Time `json:"created_at"` // CreatedAt is the timestamp indicating when the country entry was created.
	UpdatedAt time.Time `json:"updated_at"` // UpdatedAt is the timestamp indicating when the country entry was updated.
} //@name CountryDetails

// CountryTable represents a single country entry in a table.
type CountryTable struct {
	Name      string `json:"name"`       // Name is the name of the country.
	Code      string `json:"code"`       // Code is the code of the country.
	Currency  string `json:"currency"`   // Currency is the currency of the country.
	PhoneCode string `json:"phone_code"` // PhoneCode is the phone code of the country.
	Flag      string `json:"flag"`       // Flag is the flag of the country.
}
