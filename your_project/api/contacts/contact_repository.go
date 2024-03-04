package contacts

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the contact package.
type Database struct {
	DB *gorm.DB
}

// NewContactRepository performs automatic migration of contact-related structures in the database.
func NewContactRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Contact{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the contact structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of contacts based on contact ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Contact, modelID uuid.UUID, limit, offset int) ([]domains.Contact, error) {
	err := db.Where("id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a contact by its unique identifier.
func ReadContactByID(db *gorm.DB, model domains.Contact, id uuid.UUID, mailinglistID uuid.UUID) (domains.Contact, error) {
	err := db.Where("mailinglist_id= ? ", mailinglistID).First(&model, id).Error
	return model, err
}

// read all contacts for a given mailing list
func ReadAllContactsForMailingList(db *gorm.DB, mailinglistID uuid.UUID, limit, offset int) ([]domains.Contact, error) {
	var contacts []domains.Contact
	err := db.Where("mailinglist_id = ? ", mailinglistID).Limit(limit).Offset(offset).Find(&contacts).Error
	return contacts, err
}
