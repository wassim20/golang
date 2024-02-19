/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- Companies: Represents information about a Contact in the system.
		- ID (uuid.UUID): Unique identifier for the Contact.
		- Email (string): The Contact's email address.
		- Firstname (string): The firstname of the Contact.
		- Lastname (string): The lastname of the Contact.
		- PhoneNumber (string): The Contact's PhoneNumber .
		- Fullname (string): The fullname of the Contact.
		- Mailinglists ([]Mailinglist): List of Mailinglists associated with the Contact.
		- Tags ([]uuid.UUID): List of Tag's uuid associated with the Contact.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadContactNameByID(db *gorm.DB, ContactID uuid.UUID) (string, error): Retrieves the name of the Contact based on its ID from the database.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling entities in the project.

	Note:
	- The Contacts structure represents information about a Contact in the system.
	- ReadContactNameByID retrieves the name of the Contact based on its ID from the database.

	Last update :
	18/02/2024 20:26

*/

package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
	ID          uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the company
	Email       string    `gorm:"column:email; not null; unique"`              // contact's email address (unique)
	Firstname   string    `gorm:"column:first_name; not null;"`                // The contact's first name
	Lastname    string    `gorm:"column:last_name; not null;"`                 // The contact's last name
	PhoneNumber string    `gorm:"column:phone_number; not null; unique"`       // contact's email address (unique)
	FullName    string    `gorm:"column:full_name"`                            //contact's full name

	Mailinglist []Mailinglist `gorm:"many2many:mailing_list_contacts"`
	Tags        []uuid.UUID   `gorm:"column:tags;type:uuid[]"` // List of tag's uuid associated with the contact

	gorm.Model
}

func ReadContactNameByID(db *gorm.DB, ContactID uuid.UUID) (string, error) {
	contact := new(Contact)
	check := db.Select("first_name, last_name").Where("id =?", ContactID).First(contact)

	if check.Error != nil {
		return "", check.Error
	}

	return contact.Firstname + " " + contact.Lastname, nil
}
