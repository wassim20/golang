package campaign

import (
	"fmt"
	"labs/domains"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// Database represents the database instance for the companies package.
type Database struct {
	DB *gorm.DB
}

// NewCampaignRepository performs automatic migration of campaign-related structures in the database.
func NewCampaignRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Campaign{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the campaign structure. Error: ", err)
	}
}

// ReadAllPaginationFromCompany retrieves a paginated list of campaigns that belong to a specific company.
func ReadAllPaginationFromCompany(db *gorm.DB, model []domains.Campaign, companyID uuid.UUID, limit, offset int) ([]domains.Campaign, error) {
	err := db.Joins("JOIN mailinglists ON mailinglists.id = campaigns.mailinglist_id").
		Where("mailinglists.company_id = ?", companyID).
		Limit(limit).Offset(offset).
		Find(&model).Error
	return model, err
}

// ReadByID retrieves a campaign by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Campaign, id uuid.UUID) (domains.Campaign, error) {
	err := db.First(&model, id).Error
	return model, err
}

func SendCampaignEmailJob(db *gorm.DB, campaignID uuid.UUID) error {
	campaign := domains.Campaign{}
	err := db.First(&campaign, campaignID).Error
	if err != nil {
		return err
	}
	campaign.Status = "sending"
	campaign.RunAt = time.Now() // Update run time
	if err := db.Save(&campaign).Error; err != nil {
		logrus.Error("Error updating campaign status to sending:", err.Error())
		// Handle error (consider retrying or notifying admins)
	}
	mailinglist := domains.Mailinglist{}
	err = db.Preload("Contacts").First(&mailinglist, campaign.MailingListID).Error
	if err != nil {
		logrus.Errorf("can't read mailinglist from database: %v", err)
		return err
	}

	// 3. Loop Through Contacts and Send Emails
	for _, contact := range mailinglist.Contacts {
		// 3.1 Build Email Message using Gomail
		msg := gomail.NewMessage()
		msg.SetHeader("From", campaign.FromEmail)
		msg.SetHeader("Subject", campaign.Subject)
		msg.SetHeader("Reply-To", campaign.ReplyTo)
		msg.SetHeader("To", contact.Email) // Use contact email from MailingList
		// 3.2 Set Body Content (HTML or Plain Text)
		var body string
		if campaign.HTML != "" {
			body = campaign.HTML
			body = strings.Replace(body, "[Recipient Name]", contact.Firstname, -1)

		} else {
			body = campaign.Plain
		}

		//create tracking log
		trackingLog := &domains.TrackingLog{
			ID:             uuid.New(),
			CompanyID:      mailinglist.CompanyID,
			CampaignID:     campaign.ID,
			RecipientEmail: contact.Email,
			Status:         "pending",
		}

		if campaign.TrackOpen {
			// 	trackingLog.OpenTrackingID = uuid.New()
			// 	openTrackingPixelURL := "http://localhost:8080/api/" + mailinglist.CompanyID.String() + "/" + campaignID.String() + "/logs/open/" + trackingLog.OpenTrackingID.String()
			// 	// Append the tracking pixel <img> tag within the HTML body
			// 	body = strings.Replace(body, "</body>", fmt.Sprintf(`<img src="%s" width="1" height="1" alt="" style="display:none;" /></body>`, openTrackingPixelURL), 1)
		}

		if campaign.TrackClick {
			// Create a unique click ID for each link
			trackingLog.ClickTrackingID = uuid.New()

			openClickTrackingURL := "http://localhost:8080/api/" + mailinglist.CompanyID.String() + "/" + campaignID.String() + "/logs/click/" + trackingLog.ClickTrackingID.String()

			re := regexp.MustCompile(`(?i)<(a|button)[^>]*href=["'](?P<href>[^"']*)["'][^>]*>(?P<content>.*?)</(a|button)>`) // Case-insensitive match
			modifiedBody := re.ReplaceAllStringFunc(body, func(s string) string {
				matches := re.FindStringSubmatch(s)
				href := matches[re.SubexpIndex("href")]
				content := matches[re.SubexpIndex("content")]

				finalURL := href
				if href == "" {
					finalURL = "#"
				}

				// Append the tracking parameter to the original URL
				trackingURL := fmt.Sprintf(`%s?click=%s&email=%s`, finalURL, openClickTrackingURL, contact.Email)

				// Return the modified link
				return fmt.Sprintf(`<%s href="%s"%s>%s</%s>`, matches[1], trackingURL, matches[2:], content, matches[4])
			})
			body = modifiedBody

		} else {
			body = campaign.HTML
		}

		if err := domains.Create(db, trackingLog); err != nil {
			fmt.Println("***********************************************************", err)
			logrus.Error("Error saving data to the database. Error: ", err.Error())
			return err
		}
		msg.SetBody(chooseContentType(campaign.HTML, campaign.Plain), body)

		// 3.3 Optional: Add Attachments (if applicable)
		// ... (code for adding attachments based on campaign data) ...

		// 3.4 Send the Email
		d := gomail.NewDialer("smtp.gmail.com", 587, "wassimgx15@gmail.com", "zadh nbng mbdo tsbd")
		if err := d.DialAndSend(msg); err != nil {
			logrus.Error("Error sending email to", contact.Email, ":", err.Error())
			// Handle error (e.g., retry or log)
		}
	}
	// 5. Update Campaign Status to Sent
	campaign.Status = "sent"
	if err := db.Save(&campaign).Error; err != nil {
		logrus.Error("Error updating campaign status to sent:", err.Error())
		// Handle error (consider retrying or notifying admins)
	}

	return nil

}
func chooseContentType(html, plain string) string {
	if html != "" {
		return "text/html"
	}
	return "text/plain"
}
