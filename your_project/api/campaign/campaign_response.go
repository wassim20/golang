package campaign

import (
	"time"

	"github.com/google/uuid"
)

// @Description CampaignIn represents the input structure for creating a new campaign.
type CampaignIn struct {
	Type       string    `json:"type" binding:"required"`      // Type is the type of the campaign.
	Name       string    `json:"name" binding:"required"`      // Name is the name of the campaign.
	Subject    string    `json:"subject" binding:"required"`   // Subject is the subject of the campaign.
	HTML       string    `json:"html" binding:"required"`      // HTML is the HTML content of the campaign.
	FromEmail  string    `json:"fromEmail" binding:"required"` // FromEmail is the from email address of the campaign.
	FromName   string    `json:"fromName" binding:"required"`  // FromName is the from name of the campaign.
	DeliveryAt time.Time `json:"deliveryAt,omitempty"`         // DeliveryAt is the delivery time of the campaign (nullable).
	TrackOpen  bool      `json:"trackOpen"`                    // TrackOpen indicates whether opens are tracked for the campaign.
	TrackClick bool      `json:"trackClick"`                   // TrackClick indicates whether clicks are tracked for the campaign.
	ReplyTo    string    `json:"replyTo"`                      // ReplyTo is the reply-to email address of the campaign.

} //@name CampaignIn

// @Description CampaignsPagination represents the paginated list of campaigns.
type CampaignsPagination struct {
	Items      []CampaignsTable `json:"items"`      // Items is a slice containing individual campaign details.
	Page       uint             `json:"page"`       // Page is the current page number in the pagination.
	Limit      uint             `json:"limit"`      // Limit is the maximum number of items per page in the pagination.
	TotalCount uint             `json:"totalCount"` // TotalCount is the total number of campaigns in the entire list.
} //@name CampaignsPagination

// @Description CampaignsTable represents a single campaign entry in a table.
type CampaignsTable struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the campaign.
	Name        string    `json:"name"`        // Name is the name of the campaign.
	Subject     string    `json:"subject"`     // Subject is the subject of the campaign.
	FromEmail   string    `json:"fromEmail"`   // FromEmail is the from email address of the campaign.
	FromName    string    `json:"fromName"`    // FromName is the from name of the campaign.
	ReplyTo     string    `json:"replyTo"`     // ReplyTo is the reply-to email address of the campaign.
	Status      string    `json:"status"`      // Status is the status of the campaign.
	SignDKIM    bool      `json:"signDKIM"`    // SignDKIM indicates whether the campaign is signed with DKIM.
	TrackOpen   bool      `json:"trackOpen"`   // TrackOpen indicates whether opens are tracked for the campaign.
	TrackClick  bool      `json:"trackClick"`  // TrackClick indicates whether clicks are tracked for the campaign.
	Resend      bool      `json:"resend"`      // Resend indicates whether the campaign allows resend.
	CustomOrder int       `json:"customOrder"` // CustomOrder is the custom order for sorting.
	RunAt       time.Time `json:"runAt"`       // RunAt is the scheduled run time of the campaign.
	DeliveryAt  time.Time `json:"deliveryAt"`  // DeliveryAt is the delivery time of the campaign.
} //@name CampaignsTable

// @Description	CampaignsDetails represents detailed information about a specific campaign.
type CampaignsDetails struct {
	ID          uuid.UUID `json:"id"`          // ID is the unique identifier for the campaign.
	Name        string    `json:"name"`        // Name is the name of the campaign.
	Subject     string    `json:"subject"`     // Subject is the subject of the campaign.
	HTML        string    `json:"html"`        // HTML is the HTML content of the campaign.
	Plain       string    `json:"plain"`       // Plain is the plain text content of the campaign.
	FromEmail   string    `json:"fromEmail"`   // FromEmail is the from email address of the campaign.
	FromName    string    `json:"fromName"`    // FromName is the from name of the campaign.
	ReplyTo     string    `json:"replyTo"`     // ReplyTo is the reply-to email address of the campaign.
	Status      string    `json:"status"`      // Status is the status of the campaign.
	SignDKIM    bool      `json:"signDKIM"`    // SignDKIM indicates whether the campaign is signed with DKIM.
	TrackOpen   bool      `json:"trackOpen"`   // TrackOpen indicates whether opens are tracked for the campaign.
	TrackClick  bool      `json:"trackClick"`  // TrackClick indicates whether clicks are tracked for the campaign.
	Resend      bool      `json:"resend"`      // Resend indicates whether the campaign allows resend.
	CustomOrder int       `json:"customOrder"` // CustomOrder is the custom order for sorting.
	RunAt       time.Time `json:"runAt"`       // RunAt is the scheduled run time of the campaign.
	DeliveryAt  time.Time `json:"deliveryAt"`  // DeliveryAt is the delivery time of the campaign.
} //@name CampaignsDetails
