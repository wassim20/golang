package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"labs/domains"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// Database represents the database instance for the mailinglist package.
type Database struct {
	DB *gorm.DB
}

// NewWorkflowRepository performs automatic migration of workflow-related structures in the database.
func NewWorkflowRepository(db *gorm.DB) {
	if err := db.AutoMigrate(&domains.Workflow{}); err != nil {
		logrus.Fatal("An error occurred during automatic migration of the workflow structure. Error: ", err)
	}
}

// ReadAllPagination retrieves a paginated list of workflows based on mailinglist ID, limit, and offset.
func ReadAllPagination(db *gorm.DB, model []domains.Workflow, modelID uuid.UUID, limit, offset int) ([]domains.Workflow, error) {
	err := db.Where("company_id = ? ", modelID).Limit(limit).Offset(offset).Find(&model).Error
	return model, err
}

// ReadByID retrieves a workflow by its unique identifier.
func ReadByID(db *gorm.DB, model domains.Workflow, id uuid.UUID) (domains.Workflow, error) {
	err := db.Preload("Actions").First(&model, id).Error
	return model, err
}

// Start initiates the workflow execution
func Start(db *gorm.DB, workflow domains.Workflow, ctx *gin.Context, contact domains.Contact) error {
	actions, err := getOrderedActions(db, workflow.ID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	executed := make(map[uuid.UUID]bool)

	for i := 0; i < len(actions); {
		action := actions[i]
		if executed[action.ID] {
			i++
			continue
		}
		executed[action.ID] = true

		wg.Add(1)
		go func(action domains.Action) {
			defer wg.Done()
			if err := executeAction(db, workflow, action, ctx, contact); err != nil {
				logrus.Errorf("Error executing action %s: %v", action.ID, err)
			}

			// Determine next action based on condition branch if applicable
			if action.Type == "condition" {
				nextActionID := handleConditionAndDecideNext(db, workflow, action, ctx, contact)
				if nextActionID != uuid.Nil {
					// Set the loop index to the position of the next action
					for j, nextAction := range actions {
						if nextAction.ID == nextActionID {
							i = j
							break
						}
					}
				} else {
					i++ // No further actions, exit loop
				}
			} else {
				i++ // Continue to the next action
			}
		}(action)

		wg.Wait() // Wait for the action to complete before proceeding
	}

	return nil
}

// executeAction determines the type of action and handles it appropriately
func executeAction(db *gorm.DB, workflow domains.Workflow, action domains.Action, ctx *gin.Context, contact domains.Contact) error {
	switch action.Type {
	case "email":
		return handleEmailAction(db, workflow, action, contact)
	case "condition":
		// Handled separately in Start function
		return nil
	case "wait":
		return handleWaitAction(db, workflow.ID, action, ctx)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// handleEmailAction processes an email action by sending an email with tracking logs
func handleEmailAction(db *gorm.DB, workflow domains.Workflow, action domains.Action, contact domains.Contact) error {
	var emailData map[string]interface{}
	if err := json.Unmarshal([]byte(action.Data), &emailData); err != nil {
		return fmt.Errorf("error unmarshaling email data: %v", err)
	}

	requiredFields := []string{"subject", "HTML", "from", "reply_to"}
	for _, field := range requiredFields {
		if _, ok := emailData[field]; !ok {
			return fmt.Errorf("missing required field '%s' in email data", field)
		}
	}

	// Fetch available servers for sending emails
	var servers []domains.Server
	if err := db.Where("company_id = ?", workflow.CompanyID).Find(&servers).Error; err != nil {
		return fmt.Errorf("error fetching servers: %v", err)
	}
	if len(servers) == 0 {
		return fmt.Errorf("no servers available for sending emails")
	}
	server := servers[0]

	msg := gomail.NewMessage()
	msg.SetHeader("From", emailData["from"].(string))
	msg.SetHeader("To", contact.Email)
	msg.SetHeader("Subject", emailData["subject"].(string))
	msg.SetHeader("Reply-To", emailData["reply_to"].(string))

	// Create tracking log
	trackingLog := &domains.TrackingLog{
		ID:             uuid.New(),
		CompanyID:      workflow.CompanyID,
		ActionID:       action.ID,
		RecipientEmail: contact.Email,
		Status:         "pending",
	}

	body := emailData["HTML"].(string)
	body = strings.Replace(body, "[Recipient Name]", contact.Firstname, -1) // Replace placeholder with recipient's first name

	// Open tracking
	if trackOpen, ok := emailData["track_open"].(bool); ok && trackOpen {
		trackingLog.OpenTrackingID = uuid.New()
		openTrackingPixelURL := fmt.Sprintf("https://apitest385.cbot.tn/api/static/pixel.png?trackingID=%s", trackingLog.OpenTrackingID.String())
		body = strings.Replace(body, "</body>", fmt.Sprintf(`<img src="%s" width="1" height="1" alt="" style="display:none;" /></body>`, openTrackingPixelURL), 1)
	}

	// Click tracking
	if trackClick, ok := emailData["track_click"].(bool); ok && trackClick {
		trackingLog.ClickTrackingID = uuid.New()
		openClickTrackingURL := fmt.Sprintf("https://apitest385.cbot.tn/api/click?trackingID=%s", trackingLog.ClickTrackingID.String())

		// Modify links in the HTML to include click tracking
		doc, _ := html.Parse(strings.NewReader(body))

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "a" {
				for i, a := range n.Attr {
					if a.Key == "href" {
						originalURL := a.Val
						trackingURL := fmt.Sprintf("%s&redirect=%s", openClickTrackingURL, url.QueryEscape(originalURL))
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
		body = buf.String()
	}

	// Save tracking log
	if err := domains.Create(db, trackingLog); err != nil {
		logrus.Errorf("Error saving tracking log for contact %s: %v", contact.Email, err)
	}

	// Send the email
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(server.Host, server.Port, server.Username, server.Password)
	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	// Update action status to completed
	return updateActionStatus(db, action.ID, "completed")
}

// handleWaitAction processes a wait action by waiting for a specified duration
func handleWaitAction(db *gorm.DB, workflowID uuid.UUID, action domains.Action, ctx *gin.Context) error {
	var waitData map[string]interface{}
	if err := json.Unmarshal([]byte(action.Data), &waitData); err != nil {
		return fmt.Errorf("error unmarshaling wait data: %v", err)
	}

	waitDurationString, ok := waitData["duration"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'duration' in wait data")
	}
	waitDuration, err := time.ParseDuration(waitDurationString)
	if err != nil {
		return fmt.Errorf("invalid wait duration format: %v", err)
	}

	// Wait for the specified duration
	time.Sleep(waitDuration)

	// Update action status to completed
	return updateActionStatus(db, action.ID, "completed")
}

// handleConditionAndDecideNext handles a condition action and returns the ID of the next action to execute
func handleConditionAndDecideNext(db *gorm.DB, workflow domains.Workflow, action domains.Action, ctx *gin.Context, contact domains.Contact) uuid.UUID {
	var conditionData map[string]interface{}
	if err := json.Unmarshal([]byte(action.Data), &conditionData); err != nil {
		logrus.Errorf("Error unmarshaling condition data: %v", err)
		return uuid.Nil
	}

	criteria, ok := conditionData["criteria"].(string)
	if !ok {
		logrus.Errorf("Missing required field 'criteria' in condition data")
		return uuid.Nil
	}

	durationString, ok := conditionData["duration"].(string)
	if !ok {
		logrus.Errorf("Missing required field 'duration' in condition data")
		return uuid.Nil
	}
	duration, err := time.ParseDuration(durationString)
	if err != nil {
		logrus.Errorf("Invalid duration format: %v", err)
		return uuid.Nil
	}

	// Wait and check if the condition is met
	met := checkConditionWithTimeout(db, action, criteria, duration)

	// Determine the next action based on the branch outcome
	branch := "no"
	if met {
		branch = "yes"
	}
	return getNextActionByBranch(db, workflow.ID, action.ID, branch).ID
}

// checkConditionWithTimeout waits for a condition to be met within a duration
func checkConditionWithTimeout(db *gorm.DB, action domains.Action, criteria string, duration time.Duration) bool {
	timeout := time.After(duration)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return false
		case <-ticker.C:
			if checkCondition(db, action.ID, criteria) {
				return true
			}
		}
	}
}

// checkCondition checks whether the specified condition has been met
func checkCondition(db *gorm.DB, actionID uuid.UUID, criteria string) bool {
	var trackingLogs []domains.TrackingLog
	if err := db.Where("action_id = ?", actionID).Find(&trackingLogs).Error; err != nil {
		logrus.Errorf("Error fetching tracking logs: %v", err)
		return false
	}

	for _, log := range trackingLogs {
		switch criteria {
		case "read":
			if log.Status == "opened" {
				return true
			}
		case "click":
			if log.ClickCount > 0 {
				return true
			}
		}
	}

	return false
}

// getNextActionByBranch retrieves the next action based on the branch
func getNextActionByBranch(db *gorm.DB, workflowID uuid.UUID, parentID uuid.UUID, branch string) domains.Action {
	var nextAction domains.Action
	db.Where("workflow_id = ? AND parent_id = ? AND JSON_EXTRACT(data, '$.branch') = ?", workflowID, parentID, branch).First(&nextAction)
	return nextAction
}

// updateActionStatus updates the status of an action in the database
func updateActionStatus(db *gorm.DB, actionID uuid.UUID, status string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		action := domains.Action{}
		if err := tx.First(&action, actionID).Error; err != nil {
			return err
		}
		action.Status = status
		return tx.Save(&action).Error
	})
}

// getOrderedActions retrieves and orders actions based on dependencies
func getOrderedActions(db *gorm.DB, workflowID uuid.UUID) ([]domains.Action, error) {
	var actions []domains.Action
	if err := db.Where("workflow_id = ?", workflowID).Find(&actions).Error; err != nil {
		return nil, err
	}

	actionMap := make(map[uuid.UUID]*domains.Action)
	for i := range actions {
		actionMap[actions[i].ID] = &actions[i]
	}

	parentMap := make(map[uuid.UUID][]*domains.Action)
	for _, action := range actions {
		if action.ParentID != uuid.Nil {
			parentMap[action.ParentID] = append(parentMap[action.ParentID], actionMap[action.ID])
		}
	}

	var orderedActions []domains.Action
	visited := make(map[uuid.UUID]bool)
	var visit func(uuid.UUID)
	visit = func(id uuid.UUID) {
		if visited[id] {
			return
		}
		visited[id] = true
		for _, child := range parentMap[id] {
			visit(child.ID)
		}
		orderedActions = append(orderedActions, *actionMap[id])
	}

	for id := range actionMap {
		if !visited[id] {
			visit(id)
		}
	}

	return orderedActions, nil
}
