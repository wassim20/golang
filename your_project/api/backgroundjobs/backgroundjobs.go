package backgroundjobs

import (
	"bytes"
	"fmt"
	"labs/config"
	"labs/domains"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"

	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	// ... (your imports for database and campaign logic)
)

// StartBackgroundJob starts the cron job that checks for upcoming campaigns periodically.
func StartBackgroundJob(schedule string) error {

	db := config.ConnectToDB()
	crontab := cron.New()
	logrus.Info(schedule)
	logrus.Info("***************************", "starting", "***************************")

	go func() {
		crontab.AddFunc(schedule, func() {
			if err := CheckForUpcomingCampaigns(db); err != nil {
				logrus.Errorf("Error checking for upcoming campaigns: %v", err)
			}
		})
		crontab.Start()

	}()
	return nil
}

// CheckForUpcomingCampaigns checks for upcoming campaigns and triggers email sending.
func CheckForUpcomingCampaigns(db *gorm.DB) error {
	threshold := time.Now().Add(time.Second * 5)

	logrus.Debug("Checking for upcoming campaigns...")
	var campaigns []domains.Campaign
	err := db.Where("delivery_at < ? AND status = ?", threshold, "pending").Find(&campaigns).Error
	if err != nil {
		return err
	}
	if len(campaigns) == 0 {
		logrus.Debug("No upcoming campaigns found")
		return nil
	}
	fmt.Println("Found upcoming campaigns:", len(campaigns))
	for _, campaign := range campaigns {
		fmt.Printf("Upcoming campaign found: %s\n", campaign.Name)
		if err := SendCampaignEmailJob(db, campaign.ID); err != nil {
			logrus.Errorf("Error sending campaign %s: %v", campaign.Name, err)
		}
	}
	return nil
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
	// Fetch available servers for sending emails
	servers := []domains.Server{}
	err = db.Where("company_id = ?", mailinglist.CompanyID).Find(&servers).Error
	if err != nil {
		logrus.Errorf("Error fetching servers for company: %v", err)
		return err
	}
	// Check if there are any servers available
	if len(servers) == 0 {
		logrus.Error("No servers available for sending emails")
		return fmt.Errorf("no servers available for sending emails") // Or return an error if desired
	}
	emailsPerServer := len(mailinglist.Contacts) / len(servers)
	remainingContacts := len(mailinglist.Contacts) % len(servers) // Handle leftover contacts

	var wg sync.WaitGroup
	wg.Add(len(servers))

	for i, server := range servers {
		start := i * emailsPerServer
		end := (i + 1) * emailsPerServer
		if i == len(servers)-1 {
			end += remainingContacts // Assign leftovers to last server
		}

		go func(server domains.Server, start int, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				contact := mailinglist.Contacts[j]
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
					trackingLog.OpenTrackingID = uuid.New()
					//openTrackingPixelURL := "http://localhost:8080/api/" + mailinglist.CompanyID.String() + "/" + campaignID.String() + "/logs/open/" + trackingLog.OpenTrackingID.String()

					//2nd approach

					openTrackingPixelURL := fmt.Sprintf("https://apitest385.cbot.tn/api/static/pixel.png?trackingID=%s", trackingLog.OpenTrackingID.String())

					// Append the tracking pixel <img> tag within the HTML body
					body = strings.Replace(body, "</body>", fmt.Sprintf(`<img src="%s" width="1" height="1" alt="" style="display:none;" /></body>`, openTrackingPixelURL), 1)

				}

				if campaign.TrackClick {
					// Create a unique click ID for each link
					trackingLog.ClickTrackingID = uuid.New()
					openClickTrackingURL := fmt.Sprintf("https://apitest385.cbot.tn/api/click?trackingID=%s", trackingLog.ClickTrackingID.String())

					// Case-insensitive match

					// re := regexp.MustCompile(`(?i)<a\s+[^>]*href=("([^"]*)")[^>]*>(.*?)</a>`)
					// modifiedBody := re.ReplaceAllStringFunc(body, func(s string) string {
					// 	matches := re.FindStringSubmatch(s)
					// 	if len(matches) > 0 {
					// 		originalURL := matches[2]
					// 		trackingURL := fmt.Sprintf("%s?redirect=%s", openClickTrackingURL, url.QueryEscape(originalURL))
					// 		return fmt.Sprintf(`<a href="%s">%s</a>`, trackingURL, matches[3])
					// 	}
					// 	return s
					// })
					// body = modifiedBody

					//// METHOD 2 STABLE
					// doc, _ := html.Parse(strings.NewReader(body))

					// var f func(*html.Node)
					// f = func(n *html.Node) {
					// 	if n.Type == html.ElementNode && n.Data == "a" {
					// 		for i, a := range n.Attr {
					// 			if a.Key == "href" {
					// 				originalURL := a.Val
					// 				trackingURL := fmt.Sprintf("%s?redirect=%s", openClickTrackingURL, url.QueryEscape(originalURL))
					// 				n.Attr[i].Val = trackingURL
					// 			}
					// 		}
					// 	}
					// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
					// 		f(c)
					// 	}
					// }
					// f(doc)

					// var buf bytes.Buffer
					// html.Render(&buf, doc)
					// modifiedBody := buf.String()

					// body = modifiedBody

					doc, _ := html.Parse(strings.NewReader(body))

					var f func(*html.Node)
					f = func(n *html.Node) {
						if n.Type == html.ElementNode && n.Data == "a" {
							for i, a := range n.Attr {
								if a.Key == "href" {
									originalURL := a.Val
									trackingURL := fmt.Sprintf("%s?redirect=%s", openClickTrackingURL, url.QueryEscape(originalURL))
									n.Attr[i].Val = trackingURL
								}
							}
						}
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							f(c)
						}
					}
					f(doc)

					var buf bytes.Buffer
					html.Render(&buf, doc)
					modifiedBody := buf.String()

					body = modifiedBody

				} else {
					body = campaign.HTML
				}

				if err := domains.Create(db, trackingLog); err != nil {
					logrus.Errorf("error saving tracking log for contact %s: %w", contact.Email, err)
					// Handle the error (e.g., retry saving the log, log the error for debugging)
				}
				fmt.Println("Sending from server ", server.Name)
				msg.SetBody("text/html", body)
				d := gomail.NewDialer("smtp.gmail.com", 587, "wassimgx15@gmail.com", "zadh nbng mbdo tsbd")
				//d := gomail.NewDialer(server.Host, server.Port, server.Username, server.Password)
				//d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 2525, "969018208d53bf", "4da6491a9da571")
				if err := d.DialAndSend(msg); err != nil {
					logrus.Error("Error sending email to ", contact.Email, ":", err.Error())
					return
				}
				fmt.Println("/////////////////////////////////////////////////////////" + body)
			}
		}(server, start, end)
	}
	wg.Wait()

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
