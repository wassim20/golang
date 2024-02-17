package notification

import (
	"errors"
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(notification *models.Notification) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	// Validation logic using a library like "validator" (to be implemented)

	if err := tx.Create(notification).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create notification: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *NotificationService) Read(notificationID uint) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.First(&notification, notificationID).Error; err != nil {
		return nil, fmt.Errorf("Failed to read notification: %w", err)
	}
	return &notification, nil
}

func (s *NotificationService) ReadAll(preload ...string) ([]models.Notification, error) {
	var notifications []models.Notification

	if err := s.db.Find(&notifications).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No notifications found
		}
		return nil, fmt.Errorf("Failed to read notifications: %w", err)
	}

	return notifications, nil
}
func (s *NotificationService) Update(notification *models.Notification) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	// Use a library like validator to validate data before updating
	// Validate required fields and other constraints

	if err := tx.Save(notification).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update notification: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *NotificationService) Delete(notification *models.Notification) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(notification).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete notification: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}
