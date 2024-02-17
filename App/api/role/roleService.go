package role

import (
	"fmt"

	"github.com/wassim_p/App/models"
	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	if db == nil {
		panic("database connection cannot be nil")
	}
	return &RoleService{db: db}
}

func (s *RoleService) Create(role *models.Role) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Create: %v", r)
		}
	}()

	// Validation logic using a library like "validator" (to be implemented)

	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to create role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *RoleService) Read(roleID uint, preload ...string) (*models.Role, error) {
	var role models.Role

	query := s.db.Preload("Users").Preload("Users.Notifications")

	for _, rel := range preload {
		query = query.Preload(rel)
	}

	if err := query.First(&role, roleID).Error; err != nil {
		return nil, fmt.Errorf("Failed to read role: %w", err)
	}
	return &role, nil
}

func (s *RoleService) ReadAll(preload ...string) ([]models.Role, error) {
	var roles []models.Role
	query := s.db.Preload("Users").Preload("Users.Notifications")

	for _, rel := range preload {
		query = query.Preload(rel)
	}
	if err := query.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("Failed to read roles: %w", err)
	}
	return roles, nil
}

func (s *RoleService) Update(role *models.Role) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Update: %v", r)
		}
	}()

	if err := tx.Save(role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to update role: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *RoleService) Delete(role *models.Role) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Errorf("Panic in Delete: %v", r)
		}
	}()

	if err := tx.Delete(role).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("Failed to delete role: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}

	return nil
}

func (s *RoleService) AssignUser(roleID uint, userID uint) error {
	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return fmt.Errorf("Failed to find role: %w", err)
	}

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("Failed to find user: %w", err)
	}

	if err := s.db.Model(&role).Association("Users").Append(&user); err != nil {
		return fmt.Errorf("Failed to assign user to role: %w", err)
	}

	return nil
}
