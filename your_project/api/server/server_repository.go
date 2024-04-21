package server

import (
	"labs/domains"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Database represents the database instance for the server package.
type Database struct {
	DB *gorm.DB
}

// NewServerRepository performs automatic migration of server-related structures in the database.
func NewServerRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Server{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the server structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of servers based on company ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Server, modelID uuid.UUID, limit, offset int) ([]domains.Server, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a server by its unique identifier.
func ReadServerByID(db *gorm.DB, model domains.Server, id uuid.UUID) (domains.Server, error) {
	err := db.First(&model, id).Error
	return model, err
}

// ReadAllServers retrieves all servers with a limit and offset.
func ReadAllServers(db *gorm.DB, limit, offset int) ([]domains.Server, error) {
	var servers []domains.Server
	err := db.Limit(limit).Offset(offset).Find(&servers).Error
	return servers, err
}
