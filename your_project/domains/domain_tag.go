/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Companies: Represents information about a Tag in the system.
		- ID (uuid.UUID): Unique identifier for the Tag.
		- name (string): The name of the Tag.
		- CompanyID (uuid.UUID): ID of the company to which the mailinglist belongs.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadTagNameByID(db *gorm.DB, TagID uuid.UUID) (string, error): Retrieves the name of the Tag based on its ID from the database.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling entities in the project.

	Note:
	- The Tags structure represents information about a Tag in the system.
	- ReadTagNameByID retrieves the name of the Tag based on its ID from the database.

	Last update :
	18/02/2024 20:26

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the notification
	Name      string    `gorm:"column:type; not null"`                       // Type of the notification
	CompanyID uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`
	gorm.Model
}

func ReadTagNameByID(db *gorm.DB, TagID uuid.UUID) (string, error) {
	tag := new(Tag)
	check := db.Select("name").Where("id =?", TagID).First(tag)

	if check.Error != nil {
		return "", check.Error
	}

	return tag.Name, nil
}
