package user

import (
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &UserService{db: db}
}

func (s *UserService) Create(user *models.User) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	// Validation logic using a library like "validator" (to be implemented)

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create user: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *UserService) Read(userID uint, preload ...string) (*models.User, error) {
	var user models.User
	query := s.db.Preload("Notifications")
	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	if err := query.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("Failed to read user: %w", err)
	}
	return &user, nil
}

func (s *UserService) ReadAll(preload ...string) ([]models.User, error) {
	var users []models.User

	query := s.db.Preload("Notifications")
	// Add preloads for other related models based on input
	for _, rel := range preload {
		query = query.Preload(rel)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("Failed to read users: %w", err)
	}
	return users, nil
}

func (s *UserService) Update(user *models.User) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	// Use a library like validator to validate data before updating
	// Validate required fields and other constraints

	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *UserService) Delete(user *models.User) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete user: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *UserService) AssignNotification(userID uint, notificationID uint) error {
	var user models.User

	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("Failed to find user: %w", err)
	}

	var notification models.Notification
	if err := s.db.First(&notification, notificationID).Error; err != nil {
		return fmt.Errorf("Failed to find notification: %w", err)
	}

	if err := s.db.Model(&user).Association("Notifications").Append(&notification); err != nil {
		return fmt.Errorf("Failed to assign notification to user: %w", err)
	}

	return nil
}
