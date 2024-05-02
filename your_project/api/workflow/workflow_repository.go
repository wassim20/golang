package workflow

import (
	"encoding/json"
	"fmt"
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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

func Start(db *gorm.DB, workflow domains.Workflow, ctx *gin.Context) error {

	var wg sync.WaitGroup
	var x int = 2
	actions, err := getActions(db, workflow.ID)
	if err != nil {
		return err
	}
	result := make(chan bool)

	for i := 0; i < len(actions); i++ {
		action := actions[i]
		wg.Add(1)
		index := i
		go func(action domains.Action) {
			defer wg.Done()
			switch action.Type {
			case "email":
				//checking previous action status and parentID

				if x == 1 {
					x = x + 1
					return
				}
				if action.ParentID == actions[index-1].ID && index > 1 && actions[index-1].Status == "condition" {
					var checking map[string]interface{}
					x = x - 1
					if err := json.Unmarshal([]byte(action.Data), &checking); err != nil {
						logrus.Error("Error mapping request from frontend. Error: ", err.Error())
						utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
						return
					}
					// parse the data to con string
					if _, ok := checking["route"]; !ok {
						logrus.Error("Error mapping request from frontend. Error: ", err.Error())
						utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
						return
					}
					con := checking["route"].(bool)

					if con != <-result {
						return

					}
				}

				var emailData map[string]interface{}
				if err := json.Unmarshal([]byte(action.Data), &emailData); err != nil {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				// Check for required fields
				if _, ok := emailData["subject"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := emailData["track_open"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := emailData["track_click"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := emailData["HTML"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := emailData["from"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := emailData["reply-to"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if err := SendEmail(db, workflow.ID, emailData, action.ID); err != nil {
					logrus.Errorf("Error sending email action: %v", err)
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())

				}

			case "wait":
				if x == 1 {
					x = x + 1
					return
				}
				if action.ParentID == actions[index-1].ID && index > 1 && actions[index-1].Status == "condition" {
					var checking map[string]interface{}
					x = x - 1
					if err := json.Unmarshal([]byte(actions[index-1].Data), &checking); err != nil {
						logrus.Error("Error mapping request from frontend. Error: ", err.Error())
						utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
						return
					}
					// parse the data to con string
					con := checking["route"].(bool)
					if con != <-result {
						return

					}
				}
				var waitData map[string]interface{}
				if err := json.Unmarshal([]byte(action.Data), &waitData); err != nil {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := waitData["duration"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if err := Wait(db, workflow.ID, waitData, action.ID, ctx); err != nil {
					logrus.Errorf("Error waiting action: %v", err)
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
				}

			case "condition":
				// Implement condition logic here
				var conditionData map[string]interface{}
				if err := json.Unmarshal([]byte(action.Data), &conditionData); err != nil {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				// Check for required fields criteria,campaignID,duration,route
				if _, ok := conditionData["criteria"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := conditionData["campaignID"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := conditionData["duration"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if _, ok := conditionData["route"]; !ok {
					logrus.Error("Error mapping request from frontend. Error: ", err.Error())
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
					return
				}
				if err := Condition(db, workflow.ID, conditionData, action.ID, ctx, result); err != nil {
					logrus.Errorf("Error condition action: %v", err)
					utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
				}

			}

		}(action)

	}
	wg.Wait()

	return nil
}
func getActions(db *gorm.DB, workflowID uuid.UUID) ([]domains.Action, error) {
	var actions []domains.Action
	if err := db.Where("workflow_id = ?", workflowID).Find(&actions).Error; err != nil {
		return nil, err
	}

	// Build a map of actions by their ID
	actionMap := make(map[uuid.UUID]*domains.Action)
	for i := range actions {
		actionMap[actions[i].ID] = &actions[i]
	}

	// Build a map of actions by their ParentID
	parentMap := make(map[uuid.UUID][]*domains.Action)
	for _, action := range actions {
		if action.ParentID != uuid.Nil {
			parentMap[action.ParentID] = append(parentMap[action.ParentID], actionMap[action.ID])

		}
	}

	// Perform topological sort if actions form a DAG
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
		visit(id)
	}

	return orderedActions, nil
}

func SendEmail(db *gorm.DB, workflowID uuid.UUID, emailData map[string]interface{}, actionID uuid.UUID) error {
	// 1. Extract Email Details from Data
	subject, ok := emailData["subject"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'subject' in email data")
	}

	trackOpen, ok := emailData["track_open"].(bool)
	if !ok {
		return fmt.Errorf("missing required field 'track_open' in email data")
	}

	trackClick, ok := emailData["track_click"].(bool)
	if !ok {
		return fmt.Errorf("missing required field 'track_click' in email data")
	}

	htmlBody, ok := emailData["HTML"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'HTML' in email data")
	}
	from, ok := emailData["from"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'from' in email data")
	}
	replyTo, ok := emailData["reply-to"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'reply-to' in email data")
	}
	// get the workflow
	workflow := domains.Workflow{}
	err := db.First(&workflow, workflowID).Error
	if err != nil {
		return err
	}

	mailinglist := &domains.Mailinglist{}
	err = db.Preload("Contacts").First(&mailinglist, workflow.MailinglistID).Error
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
				msg.SetHeader("From", from)
				msg.SetHeader("Subject", subject)
				msg.SetHeader("Reply-To", replyTo)
				msg.SetHeader("To", contact.Email) // Use contact email from MailingList
				var body string
				if htmlBody != "" {
					body = htmlBody
					body = strings.Replace(body, "[Recipient Name]", contact.Firstname, -1)

				}

				//add the tracking here ***************IMPORTANT***************
				if trackClick && trackOpen {
					//add the tracking here
				} else {
					body = htmlBody
				}
				//create tracking log domains.Create(db, trackingLog);
				msg.SetBody("text/html", body)
				d := gomail.NewDialer(server.Host, server.Port, server.Username, server.Password)
				if err := d.DialAndSend(msg); err != nil {
					logrus.Error("Error sending email to", contact.Email, ":", err.Error())
				}

			}

		}(server, start, end)
	}

	wg.Wait()
	//update the action status
	action := domains.Action{}
	err = db.First(&action, actionID).Error
	if err != nil {
		return err
	}
	action.Status = "completed"
	if err := db.Save(&action).Error; err != nil {
		logrus.Error("Error updating action status to completed:", err.Error())
		// Handle error (consider retrying or notifying admins)
	}

	return nil
}

func Wait(db *gorm.DB, workflowID uuid.UUID, waitData map[string]interface{}, actionID uuid.UUID, ctx *gin.Context) error {
	// 1. Extract Email Details from Data
	waitDurationString, ok := waitData["duration"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'duration' (string) in wait data")
	}
	waitDuration, err := time.ParseDuration(waitDurationString)
	if err != nil {
		return fmt.Errorf("invalid wait duration format: %v", err)
	}

	// 3. Update the action status to waiting
	action := &domains.Action{}
	err = db.First(&action, actionID).Error
	if err != nil {
		return err
	}
	action.Status = "waiting"
	if err := db.Save(&action).Error; err != nil {
		logrus.Error("Error updating action status to waiting:", err.Error())
		// Handle error (consider retrying or notifying admins)
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("Error during background wait: %v", err)
				// Handle any potential errors during background execution
			}
		}()

		// Wait for the specified duration
		select {
		case <-time.After(waitDuration):
			// Action completed after waiting
			logrus.Infof("Action %s (ID: %s) waited for %s", action.Name, actionID, waitDuration)

			// Update action status to completed (consider using a transaction)
			err := db.Transaction(func(db *gorm.DB) error {
				action := &domains.Action{}
				if err := db.First(&action, actionID).Error; err != nil {
					return err
				}
				action.Status = "completed"
				if err := db.Save(&action).Error; err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				logrus.Error("Error updating action status to completed:", err.Error())
				// Handle error (consider retrying or notifying admins)
			}
		case <-ctx.Done():
			// Context canceled before wait completes
			logrus.Warnf("Action %s (ID: %s) wait canceled by context", action.Name, actionID)

			err := db.Transaction(func(tx *gorm.DB) error {
				action := &domains.Action{}
				if err := tx.First(&action, actionID).Error; err != nil {
					return err
				}
				action.Status = "canceled"
				if err := tx.Save(&action).Error; err != nil {
					return err
				}
				return nil
			})

			if err != nil {
				logrus.Error("Error updating action status to canceled:", err.Error())
				// Handle error (consider retrying or notifying admins)
			}
		}
	}()

	return nil
}
func Condition(db *gorm.DB, workflowID uuid.UUID, conditionData map[string]interface{}, actionID uuid.UUID, ctx *gin.Context, result chan bool) error {
	// 1. Extract Condition Details from Data

	criteria, ok := conditionData["criteria"].(string)
	if !ok {
		return fmt.Errorf("missing required field 'criteria' in condition data")
	}
	switch criteria {
	case "read":
		// 2. Extract Campaign ID from Data
		campaignID, ok := conditionData["campaignID"].(string)
		if !ok {
			logrus.Error("missing required field 'campaignID' in condition data")
			return fmt.Errorf("missing required field 'campaignID' in condition data")
		}
		// 3. Extract Duration from Data
		durationString, ok := conditionData["duration"].(string)
		if !ok {
			logrus.Error("missing required field 'duration' in condition data")
			return fmt.Errorf("missing required field 'duration' in condition data")
		}
		duration, err := time.ParseDuration(durationString)
		if err != nil {
			logrus.Error("invalid duration format:")
			return fmt.Errorf("invalid duration format: %v", err)
		}
		// make a go routine here
		go func() {
			//constant checking in the database while the duration is not yet reached
			ticker := time.NewTicker(time.Second * 10) // Check every second
			defer ticker.Stop()
			timer := time.NewTimer(duration)
			defer timer.Stop()

			for {
				select {
				case <-ticker.C:
					//get set of trackinglogs based on campaign id
					trackingLogs := []domains.TrackingLog{}
					err := db.Where("campaign_id = ?", campaignID).Find(&trackingLogs).Error
					if err != nil {
						logrus.Errorf("can't read trackinglogs from database: %v", err)
						return
					}
					//check if the log is read
					for _, log := range trackingLogs {
						if log.Status == "clicked" {
							// pass the log.recipient_email to the next action
							//pass the log.recipient_email to the next action

							result <- true

						}
					}
				case <-timer.C:
					// Duration has elapsed, stop the checking

					result <- false
					return
				case <-ctx.Done():
					// Context canceled before the duration is reached
					logrus.Warn("Condition check canceled by context")
					return

				}
			}

		}()
	case "click":
		campaignID, ok := conditionData["campaignID"].(string)
		if !ok {
			logrus.Error("missing required field 'campaignID' in condition data")
			return fmt.Errorf("missing required field 'campaignID' in condition data")
		}
		// 3. Extract Duration from Data
		durationString, ok := conditionData["duration"].(string)
		if !ok {
			logrus.Error("missing required field 'duration' in condition data")
			return fmt.Errorf("missing required field 'duration' in condition data")
		}
		duration, err := time.ParseDuration(durationString)
		if err != nil {
			logrus.Error("invalid duration format:")
			return fmt.Errorf("invalid duration format: %v", err)
		}
		go func() {
			//constant checking in the database while the duration is not yet reached
			ticker := time.NewTicker(time.Second * 10) // Check every second
			defer ticker.Stop()
			timer := time.NewTimer(duration)
			defer timer.Stop()

			for {
				select {
				case <-ticker.C:
					//get set of trackinglogs based on campaign id
					trackingLogs := []domains.TrackingLog{}
					err := db.Where("campaign_id = ?", campaignID).Find(&trackingLogs).Error
					if err != nil {
						logrus.Errorf("can't read trackinglogs from database: %v", err)
						return
					}
					//check if the log is read
					for _, log := range trackingLogs {
						if log.ClickCount > 0 {
							// pass the log.recipient_email to the next action

							result <- true

						}
					}
				case <-timer.C:
					// Duration has elapsed, stop the checking

					result <- false
					return
				case <-ctx.Done():
					// Context canceled before the duration is reached
					logrus.Warn("Condition check canceled by context")
					return

				}
			}

		}()

	}

	return nil
}
