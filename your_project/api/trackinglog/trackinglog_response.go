package trackinglog

import (
	"time"

	"github.com/google/uuid"
)

// @Description TrackingLogIn represents the input structure for creating a new tracking log.
type TrackingLogIn struct {
	Status string `json:"status" ` // Status indicates the current status of the tracking log.
	Error  string `json:"error"`   // Error contains any error message associated with the tracking log, if applicable.
} //@name TrackingLogIn

// @Description TrackingLogPagination represents the paginated list of tracking logs.
type TrackingLogPagination struct {
	Items      []TrackingLogTable `json:"items"`      // Items is a slice containing individual tracking log details.
	Page       uint               `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint               `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint               `json:"totalCount"` // TotalCount is the total number of tracking logs in the entire list.
} //@name TrackingLogPagination

// @Description TrackingLogTable represents a single tracking log entry in a table.
type TrackingLogTable struct {
	ID             uuid.UUID  `json:"id"`         // ID is the unique identifier for the tracking log.
	CompanyID      uuid.UUID  `json:"companyId"`  // CompanyID is the ID of the company associated with the tracking log.
	CampaignID     uuid.UUID  `json:"campaignId"` // CampaignID is the ID of the campaign associated with the tracking log.
	RecipientEmail string     `json:"recipientEmail" gorm:"column:recipient_email;"`
	OpenedAt       *time.Time `json:"openedAt" gorm:"column:opened_at;omitempty;"`
	ClickedAt      *time.Time `json:"clickedAt" gorm:"column:clicked_at;omitempty;"`
	ClickCount     int        `json:"clickCount" gorm:"column:click_count; default:0"`
	Status         string     `json:"status"`    // Status indicates the current status of the tracking log.
	Error          string     `json:"error"`     // Error contains any error message associated with the tracking log, if applicable.
	CreatedAt      time.Time  `json:"createdAt"` // CreatedAt is the timestamp indicating when the tracking log entry was created.
	UpdatedAt      time.Time  `json:"updatedAt"` // UpdatedAt is the timestamp indicating when the tracking log entry was last updated.
} //@name TrackingLogTable

// @Description TrackingLogDetails represents detailed information about a specific tracking log.
type TrackingLogDetails struct {
	CompanyID  uuid.UUID `json:"companyId"`  // CompanyID is the ID of the company associated with the tracking log.
	CampaignID uuid.UUID `json:"campaignId"` // CampaignID is the ID of the campaign associated with the tracking log.
	Status     string    `json:"status"`     // Status indicates the current status of the tracking log.
	Error      string    `json:"error"`      // Error contains any error message associated with the tracking log, if applicable.
	CreatedAt  time.Time `json:"createdAt"`  // CreatedAt is the timestamp indicating when the tracking log entry was created.
	UpdatedAt  time.Time `json:"updatedAt"`  // UpdatedAt is the timestamp indicating when the tracking log entry was last updated.
} //@name TrackingLogDetails
