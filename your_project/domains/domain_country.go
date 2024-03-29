package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// country represents the structure of a country.
type Country struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the mailinglist
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Code      string    `json:"code" gorm:"type:varchar(255);not null"`
	Currency  string    `json:"currency" gorm:"type:varchar(255);not null"`
	PhoneCode string    `json:"phone_code" gorm:"type:varchar(255);not null"`
	Flag      string    `json:"flag" gorm:"type:varchar(255);not null"`
	gorm.Model
}

// func read country name by id
func ReadCountryNameByID(db *gorm.DB, countryID uuid.UUID) (string, error) {
	country := new(Country)
	check := db.Select("name").Where("id = ?", countryID).First(country)

	if check.Error != nil {
		return "", check.Error
	}

	return country.Name, nil
}
