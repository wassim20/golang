package tag

import (
	"errors"
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type TagService struct {
	db *gorm.DB
}

func NewTagService(db *gorm.DB) *TagService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &TagService{db: db}
}

func (s *TagService) Create(tag *models.Tag) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	// Validation logic using a library like "validator" (to be implemented)

	if err := tx.Create(tag).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create tag: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TagService) Read(tagID uint) (*models.Tag, error) {
	var tag models.Tag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Tag not found
		}
		return nil, fmt.Errorf("Failed to read tag: %w", err)
	}

	return &tag, nil
}

func (s *TagService) ReadAll() ([]models.Tag, error) {
	var tags []models.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No tags found
		}
		return nil, fmt.Errorf("Failed to read tags: %w", err)
	}

	return tags, nil
}

func (s *TagService) Update(tag *models.Tag) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	if err := tx.Model(tag).Updates(tag).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update tag: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *TagService) Delete(tag *models.Tag) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(tag).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete tag: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}
