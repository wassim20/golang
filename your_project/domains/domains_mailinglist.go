/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Mailinglist: Represents information about a mailinglist in the system.
		- ID (uuid.UUID): Unique identifier for the mailinglist.
		- Name (string): The name of the mailinglist.
		- Description (string): The mailinglist's Description .
		- CompanyID (uuid.UUID): ID of the company to which the mailinglist belongs.
		- Contacts ([]Contact): List of Contacts associated with the mailinglist.
		- Tags ([]uuid.UUID): List of Tag's uuid associated with the mailinglist.
		- CreatedByUserID (uuid.UUID): ID of the user who created the mailinglist.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadMailinglistNameByID(db *gorm.DB, mailinglistID uuid.UUID) (string, error): Retrieves the name of the mailinglist based on its ID from the database.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling entities in the project.

	Note:
	- The Mailinglist structure represents information about a mailinglist in the system.
	- ReadmailinglistNameByID retrieves the name of the mailinglist based on its ID from the database.

	Last update :
	18/02/2024 20:10

*/

package domains

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Mailinglist struct {
	ID              uuid.UUID      `gorm:"column:id; primaryKey; type:uuid; not null;"`                                            // Unique identifier for the mailinglist
	Name            string         `gorm:"column:name; not null;"`                                                                 // The mailinglist's name
	Description     string         `gorm:"column:description;"`                                                                    // The mailinglist's description
	CompanyID       uuid.UUID      `gorm:"column:company_id; type:uuid; not null;"`                                                // ID of the company to which the mailinglist belongs
	Contacts        []Contact      `gorm:"foreignKey:mailinglist_id; references:id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // List of contacts associated with the mailinglist
	Tags            pq.StringArray `gorm:"column:tags;type:varchar(255)[]"`                                                        // List of tag's uuid associated with the mailinglist
	CreatedByUserID uuid.UUID      `gorm:"column:created_by_user_id; type:uuid; not null;"`                                        // ID of the user who created the mailinglist
	gorm.Model
}

func ReadMailinglistNameByID(db *gorm.DB, mailinglistID uuid.UUID) (string, error) {
	mailinglist := new(Mailinglist)
	check := db.Select("name").Where("id = ?", mailinglistID).First(mailinglist)

	if check.Error != nil {
		return "", check.Error
	}

	return mailinglist.Name, nil
}
