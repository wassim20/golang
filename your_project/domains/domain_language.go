package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// language struct represents the language entity
type Language struct {
	ID   uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the language
	Name string    `json:"name" gorm:"type:varchar(255);not null"`
	Code string    `json:"code" gorm:"type:varchar(255);not null"`
	gorm.Model
}

// func read language name by name
func ReadLanguageNameByName(db *gorm.DB, languageName string) (string, error) {
	language := new(Language)
	check := db.Select("name").Where("name = ?", languageName).First(language)

	if check.Error != nil {
		return "", check.Error
	}

	return language.Name, nil
}
