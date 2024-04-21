package domains

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// server struct represents a server to send emails from.
type Server struct {
	ID        uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"` // Unique identifier for the server
	Name      string    `gorm:"column:name;not null"`                        // Name of the server
	Host      string    `gorm:"column:host;not null"`                        // Hostname of the server
	CompanyID uuid.UUID `gorm:"column:company_id; type:uuid; not null;"`     // ID of the company to which the mailinglist belongs
	Port      int       `gorm:"column:port;not null"`                        // Port number of the server
	Type      string    `gorm:"column:type;not null"`                        // Type of server (e.g., SMTP, IMAP)
	Username  string    `gorm:"column:username;not null"`                    // Username for authentication
	Password  string    `gorm:"column:password;not null"`                    // Password for authentication
	gorm.Model
}

// read server by ID retrieves a server by its unique identifier.
func ReadServerByID(db *gorm.DB, model Server, id uuid.UUID) (Server, error) {
	err := db.First(&model, id).Error
	return model, err
}
