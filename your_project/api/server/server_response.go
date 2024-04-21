package server

import (
	"time"

	"github.com/google/uuid"
)

// @Description ServerPaginator represents the paginated list of servers.
type ServerPaginator struct {
	Items      []ServerTable `json:"items"`      // Items is a slice containing individual server details.
	Page       uint          `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint          `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint          `json:"totalCount"` // TotalCount is the total number of servers in the entire list.
} //@name ServerPaginator

// @Description ServerIn represents the input structure for creating a new server.
type ServerIn struct {
	Name     string `json:"name"`     // Name of the server.
	Host     string `json:"host"`     // Hostname of the server.
	Port     int    `json:"port"`     // Port number of the server.
	Type     string `json:"type"`     // Type of server (e.g., SMTP, IMAP).
	Username string `json:"username"` // Username for authentication.
	Password string `json:"password"` // Password for authentication.
} //@name ServerIn

// @Description ServerDetails represents detailed information about a specific server.
type ServerDetails struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the server.
	Name      string    `json:"name"`       // Name of the server.
	Host      string    `json:"host"`       // Hostname of the server.
	Port      int       `json:"port"`       // Port number of the server.
	Type      string    `json:"type"`       // Type of server (e.g., SMTP, IMAP).
	Username  string    `json:"username"`   // Username for authentication.
	Password  string    `json:"password"`   // Password for authentication.
	CreatedAt time.Time `json:"created_at"` // CreatedAt is the timestamp indicating when the server was created.
} //@name ServerDetails

// @Description ServerTable represents a single server entry in a table.
type ServerTable struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for the server.
	Name      string    `json:"name"`       // Name of the server.
	Host      string    `json:"host"`       // Hostname of the server.
	Port      int       `json:"port"`       // Port number of the server.
	Type      string    `json:"type"`       // Type of server (e.g., SMTP, IMAP).
	Username  string    `json:"username"`   // Username for authentication.
	CreatedAt time.Time `json:"created_at"` // CreatedAt is the timestamp indicating when the server was created.
} //@name ServerTable
