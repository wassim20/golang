package mailinglist

import (
	"errors"
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type MailingListService struct {
	db *gorm.DB
}

func NewMailingListService(db *gorm.DB) *MailingListService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &MailingListService{db: db}
}

func (s *MailingListService) Create(mailinglist *models.MailingList) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	if err := tx.Create(mailinglist).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create mailinglist: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MailingListService) Read(mailinglistID uint, preload ...string) (*models.MailingList, error) {
	var MailingList models.MailingList
	query := s.db.Preload("Contacts").Preload("Tags").Preload("Contacts.Tags").Preload("Contacts.MailingList")

	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.First(&MailingList, mailinglistID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // MailingList not found
		}
		return nil, fmt.Errorf("Failed to read MailingList: %w", err)
	}

	return &MailingList, nil
}

func (s *MailingListService) ReadAll(preload ...string) ([]models.MailingList, error) {
	var mailinglists []models.MailingList
	query := s.db.Preload("Contacts").Preload("Tags").Preload("Contacts.Tags").Preload("Contacts.MailingLists") // Preload MailingLists and Tags by default

	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.Find(&mailinglists).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No mailinglists found
		}
		return nil, fmt.Errorf("Failed to read mailinglists: %w", err)
	}

	return mailinglists, nil
}

func (s *MailingListService) Update(mailinglist *models.MailingList) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	if err := tx.Model(mailinglist).Updates(mailinglist).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update mailinglist: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MailingListService) Delete(mailinglist *models.MailingList) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(mailinglist).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete mailinglist: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *MailingListService) AssignContact(mailingListID uint, contactID uint) error {
	var mailingList models.MailingList
	if err := s.db.First(&mailingList, mailingListID).Error; err != nil {
		return fmt.Errorf("Failed to find mailing list: %w", err)
	}

	var contact models.Contact
	if err := s.db.First(&contact, contactID).Error; err != nil {
		return fmt.Errorf("Failed to find contact: %w", err)
	}

	if err := s.db.Model(&mailingList).Association("Contacts").Append(&contact); err != nil {
		return fmt.Errorf("Failed to assign contact to mailing list: %w", err)
	}

	return nil
}

func (s *MailingListService) AssignTag(mailingListID uint, tagID uint) error {
	var mailingList models.MailingList
	if err := s.db.First(&mailingList, mailingListID).Error; err != nil {
		return fmt.Errorf("Failed to find mailing list: %w", err)
	}

	var tag models.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		return fmt.Errorf("Failed to find tag: %w", err)
	}

	//nejbdou el company bech nchoufou el tag esque feha wala le
	var company models.Company
	if err := s.db.First(&company, mailingList.CompanyID).Error; err != nil {
		return fmt.Errorf("Failed to find company for mailing list: %w", err)
	}

	fmt.Println("heeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeere")

	if company.ID != tag.CompanyID {
		return errors.New("Tag is not associated with the company")
	}

	if err := s.db.Model(&mailingList).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("Failed to assign tag to mailing list: %w", err)
	}

	return nil
}
