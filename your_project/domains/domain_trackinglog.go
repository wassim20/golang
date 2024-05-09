/*

	Package domains provides the data structures representing entities in the project.

	Structures:
	- TrackingLog: Represents information about a tracking in the system.
		- ID (uuid.UUID): Unique identifier for the Tag.
		- ContactID (uuid.UUID): ID of the contact to which the tracking belongs.
		- CampaignID (uuid.UUID): ID of the campaign to which the tracking belongs.
		- Status (string): Status of the tracking.
		- Error (string): Error message of the tracking.
		- gorm.Model: Standard GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).

	Functions:
	- ReadTrackStatusByID(db *gorm.DB, TagID uuid.UUID) (string, error): Retrieves the status of the tracking based on its ID from the database.

	Dependencies:
	- "github.com/google/uuid": Package for working with UUIDs.
	- "gorm.io/gorm": The GORM library for object-relational mapping in Go.

	Usage:
	- Import this package to utilize the provided data structures and functions for handling entities in the project.

	Note:
	- The TrackingLog structure represents information about a tracking in the system.
	- ReadTrackStatusByID retrieves the status of the tracking based on its ID from the database.

	Last update :
	12/03/2024 01:50

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrackingLog struct {
	ID              uuid.UUID  `gorm:"column:id; primaryKey; type:uuid; not null;"`
	CompanyID       uuid.UUID  `gorm:"column:company_id; type:uuid; not null;"`
	CampaignID      uuid.UUID  `gorm:"column:campaign_id; type:uuid; not null;"`
	ActionID        uuid.UUID  `gorm:"column:action_id; type:uuid; not null;"`
	Status          string     `gorm:"column:status; not null"`
	Error           string     `gorm:"nullable"`
	RecipientEmail  string     `gorm:"column:recipient_email; "`      // Email address of the recipient
	OpenTrackingID  uuid.UUID  `gorm:"column:open_tracking_id;"`      // Unique ID for open tracking (if used)
	OpenedAt        time.Time  `gorm:"column:opened_at;"`             // Timestamp when the email was opened (optional)
	ClickTrackingID uuid.UUID  `gorm:"column:click_tracking_id;"`     // Unique ID for click tracking (if used)
	ClickedAt       *time.Time `gorm:"column:clicked_at;omitempty;"`  // Timestamp when a link was clicked (optional)
	ClickCount      int        `gorm:"column:click_count; default:0"` // Number of times a link was clicked
	gorm.Model
}

func ReadTrackStatusByID(db *gorm.DB, ID uuid.UUID) (string, error) {
	log := new(TrackingLog)
	check := db.Select("status").Where("id =?", ID).First(log)

	if check.Error != nil {
		return "", check.Error
	}

	return log.Status, nil
}
