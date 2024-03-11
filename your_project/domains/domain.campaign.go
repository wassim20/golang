/*

	// Package domains provides the data structures representing entities in the project.

// Structures:
//  - Campaign: Represents information about a campaign in the system.
//     		Unique identifiers (ID, CreatedByUserID, MailingListID).
// 			Campaign details (Type, Name, Subject, HTML, Plain).
//			Sender information (FromEmail, FromName, ReplyTo).
//			Campaign status (Status).
//			Email tracking settings (SignDKIM, TrackOpen, TrackClick).
//			Resend and sorting options (Resend, CustomOrder).
//			Scheduled and delivery times (RunAt, DeliveryAt).
// Functions:
//  - ReadCampaignNameByID(db *gorm.DB, campaignID uuid.UUID) (string, error): Retrieves the name of the campaign based on its ID from the database.

// Dependencies:
//  - "github.com/google/uuid": Package for working with UUIDs.
//  - "gorm.io/gorm": The GORM library for object-relational mapping in Go.

// Usage:
//  - Import this package to utilize the provided data structures and functions for handling entities in the project.

// Note:
//  - The Campaign structure represents information about a campaign in the system.
//  - ReadCampaignNameByID retrieves the name of the campaign based on its ID from the database.

// Last update:
//  2024-03-04 12:00

*/

package domains

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Campaign struct {
	ID              uuid.UUID `gorm:"column:id; primaryKey; type:uuid; not null;"`     // Unique identifier for the campaign
	CreatedByUserID uuid.UUID `gorm:"column:created_by_user_id; type:uuid; not null;"` // ID of the user who created the campaign
	MailingListID   uuid.UUID `gorm:"column:mailinglist_id; not null;"`                // Foreign Key referencing the mail_lists table
	Type            string    `gorm:"column:type; type:varchar(255)"`                  // Campaign type
	Name            string    `gorm:"column:name; type:varchar(255)"`                  // Campaign name
	Subject         string    `gorm:"column:subject; type:text"`                       // Campaign subject
	HTML            string    `gorm:"column:html; type:text"`                          // Campaign HTML content
	Plain           string    `gorm:"column:plain; type:text"`                         // Campaign plain text content
	FromEmail       string    `gorm:"column:from_email; type:varchar(255)"`            // From email address
	FromName        string    `gorm:"column:from_name; type:varchar(255)"`             // From name
	ReplyTo         string    `gorm:"column:reply_to; type:varchar(255)"`              // Reply-to email address
	Status          string    `gorm:"column:status; type:varchar(255)"`                // Campaign status
	SignDKIM        bool      `gorm:"column:sign_dkim; default:false"`                 // Sign with DKIM DomainKeys Identified Mail
	TrackOpen       bool      `gorm:"column:track_open; default:false"`                // Track opens
	TrackClick      bool      `gorm:"column:track_click; default:false"`               // Track clicks
	Resend          bool      `gorm:"column:resend; default:false"`                    // Allow resend
	CustomOrder     int       `gorm:"column:custom_order; default:0"`                  // Custom order for sorting
	RunAt           time.Time `gorm:"column:run_at; type:timestamp"`                   // Scheduled run time (nullable)
	DeliveryAt      time.Time `gorm:"column:delivery_at; type:timestamp"`              // Delivery time (nullable)
	gorm.Model
}

func ReadCampaignNameByID(db *gorm.DB, campaignID uuid.UUID) (string, error) {
	campaign := new(Campaign)
	check := db.Select("name").Where("id = ?", campaignID).First(campaign)

	if check.Error != nil {
		return "", check.Error
	}

	return campaign.Name, nil
}
