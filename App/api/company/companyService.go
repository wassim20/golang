package company

import (
	"errors"
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type CompanyService struct {
	db *gorm.DB
}

func NewCompanyService(db *gorm.DB) *CompanyService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &CompanyService{db: db}
}

func (s *CompanyService) Create(company *models.Company) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	// Validation logic using a library like "validator" (to be implemented)

	if err := tx.Create(company).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create company: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CompanyService) Read(companyID uint, preload ...string) (*models.Company, error) {
	var company models.Company
	query := s.db.Preload("MailingLists").Preload("Tags").Preload("MailingLists.Contacts").Preload("MailingLists.Tags").Preload("Users").Preload("MailingLists.Contacts.Tags").Preload("Users.Notifications") // Preload MailingLists by default

	// Add preloads for other related models based on input ("Tags", etc.)
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.First(&company, companyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Company not found
		}
		return nil, fmt.Errorf("Failed to read company: %w", err)
	}

	return &company, nil
}

func (s *CompanyService) ReadAll(preload ...string) ([]models.Company, error) {

	var companies []models.Company
	query := s.db.Preload("MailingLists").Preload("Tags").Preload("MailingLists.Contacts").Preload("MailingLists.Tags").Preload("Users").Preload("MailingLists.Contacts.Tags").Preload("Users.Notifications") // Preload MailingLists and Tags by default

	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.Find(&companies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No companies found
		}
		return nil, fmt.Errorf("Failed to read companies: %w", err)
	}

	return companies, nil
}

func (s *CompanyService) Update(company *models.Company) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	// Use a library like validator to validate data before updating
	// Validate required fields and other constraints

	if err := tx.Model(company).Updates(company).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update company: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CompanyService) Delete(company *models.Company) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(company).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete company: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CompanyService) AssignMailingList(companyID uint, mailingListID uint) error {
	var company models.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		return fmt.Errorf("Failed to find company: %w", err)
	}

	var mailingList models.MailingList
	if err := s.db.First(&mailingList, mailingListID).Error; err != nil {
		return fmt.Errorf("Failed to find mailing list: %w", err)
	}

	if err := s.db.Model(&company).Association("MailingLists").Append(&mailingList); err != nil {
		return fmt.Errorf("Failed to assign mailing list to company: %w", err)
	}

	return nil
}

func (s *CompanyService) AssignTag(companyID uint, tagID uint) error {
	var company models.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		return fmt.Errorf("Failed to find company: %w", err)
	}

	var tag models.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		return fmt.Errorf("Failed to find tag: %w", err)
	}

	if err := s.db.Model(&company).Association("Tags").Append(&tag); err != nil {
		return fmt.Errorf("Failed to assign tag to company: %w", err)
	}

	return nil
}

func (s *CompanyService) AssignUser(companyID uint, userID uint) error {

	var company models.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		return fmt.Errorf("Failed to find company: %w", err)
	}
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("Failed to find user: %w", err)
	}

	if err := s.db.Model(&company).Association("Users").Append(&user); err != nil {
		return fmt.Errorf("Failed to assign user to company: %w", err)
	}
	return nil

}
