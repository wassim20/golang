package contact

import (
	"errors"
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &ContactService{db: db}
}

func (s *ContactService) Create(contact *models.Contact) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	if err := tx.Create(contact).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create contact: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ContactService) Read(contactID uint, preload ...string) (*models.Contact, error) {
	var contact models.Contact

	query := s.db.Preload("MailingLists").Preload("Tags")
	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	if err := query.First(&contact, contactID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Contact not found
		}
		return nil, fmt.Errorf("Failed to read contact: %w", err)
	}

	return &contact, nil
}

func (s *ContactService) ReadAll(preload ...string) ([]models.Contact, error) {
	var contacts []models.Contact

	query := s.db.Preload("MailingLists").Preload("Tags")
	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.Find(&contacts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No contacts found
		}
		return nil, fmt.Errorf("Failed to read contacts: %w", err)
	}

	return contacts, nil
}

func (s *ContactService) Update(contact *models.Contact) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	if err := tx.Model(contact).Updates(contact).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update contact: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ContactService) Delete(contact *models.Contact) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(contact).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete contact: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *ContactService) AssignTag(contactID uint, tagID uint) error {
	var contact models.Contact
	if err := s.db.First(&contact, contactID).Error; err != nil {
		return fmt.Errorf("Failed to find contact: %w", err)
	}

	var tag models.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		return fmt.Errorf("Failed to find tag: %w", err)
	}

	// same lezemna nal9aw el mailinglist w mel mailing list lezm nal9aw el company bech nchoufou el tag fiha wala le
	var mailingLists []models.MailingList
	if err := s.db.Model(&contact).Association("MailingLists").Find(&mailingLists); err != nil {
		return fmt.Errorf("Failed to find mailing lists for contact: %w", err)
	}
	fmt.Printf("heeeeeeeeeeeeeeeeeere")
	fmt.Println(mailingLists[0].CompanyID)

	if len(mailingLists) == 0 {
		return errors.New("Contact has no associated mailing lists")
	}

	// Check if the tag is associated with the company of the first mailing list

	if tag.CompanyID != mailingLists[0].CompanyID {
		return errors.New("Tag is not associated with the company")
	}

	if err := s.db.Model(&contact).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("Failed to assign tag to contact: %w", err)
	}

	return nil
}
