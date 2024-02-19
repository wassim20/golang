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
func ReadContactByID(db *gorm.DB, model domains.Contact, id uuid.UUID) (domains.Contact, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllContacts retrieves all contacts associated with a specific mailinglist ID.
func ReadAllContactsForMailingList(db *gorm.DB, mailingListID uuid.UUID) ([]domains.Contact, error) {
	var mailingList domains.Mailinglist
	var mailingListContacts []domains.Contact

	// Find the mailing list by its ID
	if err := db.First(&mailingList, mailingListID).Error; err != nil {
		return nil, err
	}

	// Retrieve all contacts associated with the mailing list
	if err := db.Model(&mailingList).Association("Contacts").Find(&mailingListContacts); err != nil {
		return nil, err
	}

	return mailingListContacts, nil
}

func ReadAllContactsForAllMailingLists(db *gorm.DB, model []domains.Contact, CompanyID uuid.UUID, limit, offset int) ([]domains.Contact, error) {
	var mailingLists []domains.Mailinglist
	if err := db.Where("company_id = ?", CompanyID).Find(&mailingLists).Error; err != nil {
		return nil, err
	}

	// Iterate over each mailing list to retrieve its contacts
	for _, mailingList := range mailingLists {
		var mailingListContacts []domains.Contact
		if err := db.Model(&mailingList).Association("Contacts").Find(&mailingListContacts); err != nil {
			return nil, err
		}

		// Append mailing list contacts to the contacts slice
		model = append(model, mailingListContacts...)
	}

	return model, nil
}
