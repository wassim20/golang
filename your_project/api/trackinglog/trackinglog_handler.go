package trackinglog

import (
	"fmt"
	"labs/api/campaign"
	"labs/constants"
	"labs/domains"
	"net/http"
	"sort"
	"strconv"
	"time"

	"labs/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateTrackingLog handles the creation of a new tracking log.
// @Summary      Create tracking log
// @Description  Create a new tracking log.
// @Tags         TrackingLogs
// @Accept        json
// @Produce      json
// @Param        request body TrackingLogIn true "Tracking log data"
// @Success      201     {object} utils.ApiResponses "Tracking log created successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/{campaignID}/logs [post]
func (db Database) CreateTrackingLog(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("campaignID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a TrackingLogIn struct
	log := new(TrackingLogIn)
	if err := ctx.ShouldBindJSON(log); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	var campaigndb domains.Campaign
	campaigndb, err = campaign.ReadByID(db.DB, campaigndb, campaignID)
	if err != nil {
		logrus.Error("Error getting campaign by ID. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Create a new user in the database
	dbLog := &domains.TrackingLog{
		ID:         uuid.New(),
		CompanyID:  companyID,
		CampaignID: campaignID,
		Status:     "pending",
		Error:      "",
		ClickedAt:  nil,
	}

	if campaigndb.TrackOpen {
		dbLog.OpenTrackingID = uuid.New()
	}

	if campaigndb.TrackClick {
		dbLog.ClickTrackingID = uuid.New()
	}

	if err := domains.Create(db.DB, dbLog); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadTrackingLogs handles the retrieval of all or paginated tracking logs.
// @Summary      Get tracking logs
// @Description  Get all tracking logs or paginated results.
// @Tags         TrackingLogs
// @Produce      json
// @Param        page query int false "Page number (defaults to 1)"
// @Param        limit query int false "Limit per page (defaults to 10)"
// @Success      200     {object} TrackingLogPagination "List of tracking logs"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/{campaignID}/logs [get]
func (db Database) ReadTrackingLogs(ctx *gin.Context) {

	// Extract JWT values from the context
	//session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("campaignID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	// if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
	// 	logrus.Error("Error verifying employee belonging. Error: ", err.Error())
	// 	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
	// 	return
	// }

	// Parse and validate the page from the request parameter
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(constants.DEFAULT_PAGE_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the limit from the request parameter
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(constants.DEFAULT_LIMIT_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the user's value is among the allowed choices
	validChoices := utils.ResponseLimitPagination()
	isValidChoice := false
	for _, choice := range validChoices {
		if uint(limit) == choice {
			isValidChoice = true
			break
		}
	}

	// If the value is invalid, set it to default DEFAULT_LIMIT_PAGINATION
	if !isValidChoice {
		limit = constants.DEFAULT_LIMIT_PAGINATION
	}

	// Generate offset
	offset := (page - 1) * limit

	// Retrieve all company data from the database
	logs, err := ReadAllPagination(db.DB, []domains.TrackingLog{}, companyID, campaignID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := ReadTotalCountTrackingLog(db.DB, companyID, campaignID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if len(logs) == 0 {
		// No logs found, return a specific response
		utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
			"message": "No tracking logs available for the specified company and campaign.",
		})
		return
	}

	// Generate a log structure as a response
	response := TrackingLogPagination{}
	listlogs := []TrackingLogTable{}
	for _, log := range logs {

		listlogs = append(listlogs, TrackingLogTable{
			ID:             log.ID,
			CompanyID:      log.CompanyID,
			CampaignID:     log.CampaignID,
			Status:         log.Status,
			Error:          log.Error,
			RecipientEmail: log.RecipientEmail,
			OpenedAt:       log.OpenedAt,
			ClickedAt:      log.ClickedAt,
			ClickCount:     log.ClickCount,
			CreatedAt:      log.CreatedAt,
			UpdatedAt:      log.UpdatedAt,
		})
	}
	response.Items = listlogs
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = uint(count)

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

func (db Database) ReadAllTrackingLogs(ctx *gin.Context) {

	// Extract JWT values from the context
	//session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	// if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
	// 	logrus.Error("Error verifying employee belonging. Error: ", err.Error())
	// 	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
	// 	return
	// }

	// Parse and validate the page from the request parameter
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(constants.DEFAULT_PAGE_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the limit from the request parameter
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", strconv.Itoa(constants.DEFAULT_LIMIT_PAGINATION)))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid INT format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the user's value is among the allowed choices
	validChoices := utils.ResponseLimitPagination()
	isValidChoice := false
	for _, choice := range validChoices {
		if uint(limit) == choice {
			isValidChoice = true
			break
		}
	}

	// If the value is invalid, set it to default DEFAULT_LIMIT_PAGINATION
	if !isValidChoice {
		limit = constants.DEFAULT_LIMIT_PAGINATION
	}

	// Generate offset
	offset := (page - 1) * limit

	// Retrieve all company data from the database
	logs, err := ReadAll(db.DB, []domains.TrackingLog{}, companyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := ReadTotalCountAllTrackingLog(db.DB, companyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if len(logs) == 0 {
		// No logs found, return a specific response
		utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, gin.H{
			"message": "No tracking logs available for the specified company and campaign.",
		})
		return
	}

	// Generate a log structure as a response
	response := TrackingLogPagination{}
	listlogs := []TrackingLogTable{}
	for _, log := range logs {

		listlogs = append(listlogs, TrackingLogTable{
			ID:             log.ID,
			CompanyID:      log.CompanyID,
			CampaignID:     log.CampaignID,
			Status:         log.Status,
			Error:          log.Error,
			RecipientEmail: log.RecipientEmail,
			OpenedAt:       log.OpenedAt,
			ClickedAt:      log.ClickedAt,
			ClickCount:     log.ClickCount,
			CreatedAt:      log.CreatedAt,
			UpdatedAt:      log.UpdatedAt,
		})
	}
	response.Items = listlogs
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = uint(count)

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadTrackingLogByID handles the retrieval of a tracking log by ID.
// @Summary      Get tracking log by ID
// @Description  Get a specific tracking log by its ID.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      200     {object} TrackingLog "Tracking log details"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router			 /{companyID}/{campaignID}/logs/{ID}	[get]
func (db Database) ReadTrackingLogByID(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("campaignID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the company ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the log data by ID from the database
	log, err := ReadByID(db.DB, domains.TrackingLog{}, companyID, campaignID, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a Log structure as a response
	details := TrackingLogDetails{
		CompanyID:  log.CompanyID,
		CampaignID: log.CampaignID,
		Status:     log.Status,
		Error:      log.Error,
		CreatedAt:  log.CreatedAt,
		UpdatedAt:  log.UpdatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateTrackingLog handles the update of a tracking log.
// @Summary      Update tracking log
// @Description  Update an existing tracking log.
// @Tags         TrackingLogs
// @Accept        json
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Param        request body TrackingLogIn true "Tracking log update data"
// @Success      200     {object} utils.ApiResponses "Tracking log updated successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router			/{companyID}/{campaignID}/logs/{ID}	[put]
func (db Database) UpdateTrackingLog(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("campaignID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the log ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a TrackingLogIn struct
	log := new(TrackingLogIn)
	if err := ctx.ShouldBindJSON(log); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the log with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.TrackingLog{}, objectID); err != nil {
		logrus.Error("Error checking if the log with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the log data in the database
	dblog := &domains.TrackingLog{
		Status:     log.Status,
		CampaignID: campaignID,
		Error:      log.Error,
	}
	if err := domains.Update(db.DB, dblog, objectID); err != nil {
		logrus.Error("Error updating company data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteTrackingLog handles the deletion of a tracking log.
// @Summary      Delete tracking log
// @Description  Delete a tracking log by ID.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      204     "Tracking log deleted successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/{campaignID}/logs/{ID} [delete]
func (db Database) DeleteTrackingLog(ctx *gin.Context) {

	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// // Parse and validate the log ID from the request parameter
	// campaignID, err := uuid.Parse(ctx.Param("campaignID"))
	// if err != nil {
	// 	logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
	// 	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
	// 	return
	// }
	// Parse and validate the log ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified log
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the log with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.TrackingLog{}, objectID); err != nil {
		logrus.Error("Error checking if the log with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Delete the log data from the database
	if err := domains.Delete(db.DB, &domains.TrackingLog{}, objectID); err != nil {
		logrus.Error("Error deleting log data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// handleOpenRequest handles the opening of an email and updates tracking log.
// @Summary      updates tracking log on email open
// @Description  updates a tracking log by ID when email open.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      204     "Tracking log deleted successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/{campaignID}/logs/open/{trackingID} [delete]
func (db Database) handleOpenRequest(ctx *gin.Context) {
	trackingID := ctx.Query("trackingID")
	if trackingID == "" {
		logrus.Error("Missing trackingID in query parameters")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var trackingLog domains.TrackingLog
	err := db.DB.First(&trackingLog, "open_tracking_id = ?", trackingID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Tracking ID not found, handle error (e.g., log or return appropriate status code)
			logrus.Errorf("Open tracking ID '%s' not found", trackingID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error fetching tracking log for open tracking ID '%s': %v", trackingID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Update tracking log status (optional)
	trackingLog.Status = "opened" // Update status if needed
	openedAt := time.Now()
	trackingLog.OpenedAt = &openedAt // Update opened_at timestamp if needed
	trackingLog.ClickCount++         // Increment click count if needed

	if err := db.DB.Save(&trackingLog).Error; err != nil {
		logrus.Errorf("Error updating tracking log status for open tracking ID '%s': %v", trackingID, err)
		// Handle error (consider logging or retrying update)
	}

	// Respond to the request (optional)
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// handleClickRequest handles the clicking of an email and updates tracking log.
// @Summary      updates tracking log on email click
// @Description  updates a tracking log by ID when email click.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      204     "Tracking log deleted successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/{campaignID}/logs/click/{trackingID} [delete]
func (db Database) handleClickRequest(ctx *gin.Context) {
	trackingID := ctx.Param("trackingID")
	if trackingID == "" {
		logrus.Error("Tracking ID is required")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update the tracking log with the click information
	trackingLog := &domains.TrackingLog{}
	if err := db.DB.First(trackingLog, "click_tracking_id = ?", trackingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Tracking ID not found, handle error (e.g., log or return appropriate status code)
			logrus.Errorf("Click tracking ID '%s' not found", trackingID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error fetching tracking log for click tracking ID '%s': %v", trackingID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Update the tracking log with click information
	clickedAt := time.Now()
	trackingLog.Status = "clicked"
	trackingLog.ClickedAt = &clickedAt
	trackingLog.ClickCount++

	if err := db.DB.Save(trackingLog).Error; err != nil {
		logrus.Errorf("Error updating tracking log: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Respond to the request (optional)
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// handleOpenRequestWorflow handles the clicking of an email and updates tracking log.
// @Summary      updates tracking log on email open
// @Description  updates a tracking log by ID when email open.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      204     "Tracking log updated successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/logs/open/{trackingID} [update]
// handle read request comming from workflow "http://localhost:8080/api/" + workflow.CompanyID.String() + "/logs/open/" + trackingLog.OpenTrackingID.String()
func (db Database) handleOpenRequestWorflow(ctx *gin.Context) {
	trackingID := ctx.Param("trackingID")

	var trackingLog domains.TrackingLog
	err := db.DB.First(&trackingLog, "open_tracking_id = ?", trackingID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Tracking ID not found, handle error (e.g., log or return appropriate status code)
			logrus.Errorf("Open tracking ID '%s' not found", trackingID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error fetching tracking log for open tracking ID '%s': %v", trackingID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Update tracking log status (optional)
	trackingLog.Status = "opened" // Update status if needed
	openedAt := time.Now()
	trackingLog.OpenedAt = &openedAt // Update opened_at timestamp if needed
	// Increment click count if needed

	if err := db.DB.Save(&trackingLog).Error; err != nil {
		logrus.Errorf("Error updating tracking log status for open tracking ID '%s': %v", trackingID, err)
		// Handle error (consider logging or retrying update)
	}

	// Respond to the request (optional)
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// handleClickRequestWorflow handles the clicking of an email and updates tracking log.
// @Summary      updates tracking log on email click
// @Description  updates a tracking log by ID when email click.
// @Tags         TrackingLogs
// @Produce      json
// @Param        id path string true "Tracking log ID"
// @Success      204     "Tracking log updated successfully"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      404     {object} utils.ApiResponses "Tracking log not found"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/logs/click/{trackingID} [update]
func (db Database) handleClickRequestWorflow(ctx *gin.Context) {
	trackingID := ctx.Param("trackingID")
	if trackingID == "" {
		logrus.Error("Tracking ID is required")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Update the tracking log with the click information
	trackingLog := &domains.TrackingLog{}
	if err := db.DB.First(trackingLog, "click_tracking_id = ?", trackingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Tracking ID not found, handle error (e.g., log or return appropriate status code)
			logrus.Errorf("Click tracking ID '%s' not found", trackingID)
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		logrus.Errorf("Error fetching tracking log for click tracking ID '%s': %v", trackingID, err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Update the tracking log with click information
	clickedAt := time.Now()
	trackingLog.Status = "clicked"
	trackingLog.ClickedAt = &clickedAt
	trackingLog.ClickCount++

	if err := db.DB.Save(trackingLog).Error; err != nil {
		logrus.Errorf("Error updating tracking log: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Respond to the request (optional)
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// updateChartData handles the count and sending of barchart data.
// @Summary      Get bar chart data
// @Description  Get the count of opened, clicked and error tracking logs of a specific campaign.
// @Tags         TrackingLogs
// @Produce      json
// @Param        companyID path string
// @Param        campaignID path string
// @Success      200     {object} utils.ApiResponses "Bar chart data"
// @Failure      400     {object} utils.ApiResponses "Invalid request"
// @Failure      401     {object} utils.ApiResponses "Unauthorized"
// @Failure      403     {object} utils.ApiResponses "Forbidden"
// @Failure      500     {object} utils.ApiResponses "Internal Server Error"
// @Router       /{companyID}/logs/barchartdata [get]
func (db Database) updateChartData(ctx *gin.Context) {
	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var result struct {
		OpenedCount  int64 `json:"opened"`
		ClickedCount int64 `json:"clicked"`
		ErrorCount   int64 `json:"error"`
	}

	// Make a query to get the count of opened, clicked, and error tracking logs of a specific campaign
	if campaignID == "" {

		err := db.DB.Model(&domains.TrackingLog{}).
			Select("COUNT(CASE WHEN opened_at IS NOT NULL AND opened_at  != '0001-01-01 00:09:21+00:09:21' THEN 1 ELSE NULL END) AS opened_count, "+
				"COUNT(CASE WHEN clicked_at IS NOT NULL THEN 1 ELSE NULL END) AS clicked_count, "+
				"COUNT(CASE WHEN error IS NOT NULL AND error !='' THEN 1 ELSE NULL END) AS error_count").
			Where("company_id = ?", companyID).
			Scan(&result).Error

		if err != nil {
			logrus.Error("Error occurred while finding all company data. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	} else {
		err := db.DB.Model(&domains.TrackingLog{}).
			Select("COUNT(CASE WHEN opened_at IS NOT NULL AND opened_at  != '0001-01-01 00:09:21+00:09:21' THEN 1 ELSE NULL END) AS opened_count, "+
				"COUNT(CASE WHEN clicked_at IS NOT NULL THEN 1 ELSE NULL END) AS clicked_count, "+
				"COUNT(CASE WHEN error IS NOT NULL AND error !='' THEN 1 ELSE NULL END) AS error_count").
			Where("company_id = ? AND campaign_id = ?", companyID, campaignID).
			Scan(&result).Error

		if err != nil {
			logrus.Error("Error occurred while finding all company data. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, result)
}
func (db Database) updatePieChartData(ctx *gin.Context) {

	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var logs []domains.TrackingLog

	// Fetch logs
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	openedEmails := 0
	clickedEmails := 0

	for _, log := range logs {

		if log.OpenedAt != nil {
			openedEmails++
		}
		if log.ClickedAt != nil {
			clickedEmails++
		}
	}
	var result struct {
		OpenedEmails  int64 `json:"cpenedEmails"`
		ClickedEmails int64 `json:"clickedEmails"`
	}
	result.OpenedEmails = int64(openedEmails)
	result.ClickedEmails = int64(clickedEmails)

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, result)
}
func (db Database) updateRadialChartData(ctx *gin.Context) {
	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var logs []domains.TrackingLog

	// Fetch logs
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	totalLogs := len(logs)
	openedLogs := 0

	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			openedLogs++
		}
	}
	openedPercentage := float64(0)
	if totalLogs > 0 {
		openedPercentage = float64(openedLogs) / float64(totalLogs) * 100
	}

	type Response struct {
		OpenedPercentage float64 `json:"openedPercentage"`
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, Response{OpenedPercentage: openedPercentage})
}
func aggregateDataByDate(logs []domains.TrackingLog, key string) map[string]int {
	if logs == nil {
		logrus.Error("Logs data is null or undefined")
		return map[string]int{}
	}
	aggregatedData := make(map[string]int)
	for _, log := range logs {
		var timestamp *time.Time
		switch key {
		case "openedAt":
			timestamp = log.OpenedAt
		case "clickedAt":
			timestamp = log.ClickedAt
		}
		if !timestamp.IsZero() && timestamp.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			date := timestamp.Format("2006-01-02")
			aggregatedData[date]++
		}
	}
	return aggregatedData
}
func (db Database) updateLineChartData(ctx *gin.Context) {

	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var openedData []domains.TrackingLog
	var clickedData []domains.TrackingLog
	var logs []domains.TrackingLog

	if campaignID == "" {
		// Fetch logs and pass opened and clicked logs to the function aggregate simultaniously

		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

	} else {
		// Fetch logs and pass opened and clicked logs to the function aggregate simultaniously

		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}
	// Separate logs into openedData and clickedData
	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			openedData = append(openedData, log)
		}
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" { // Check if ClickedAt is not nil before using it
			clickedData = append(clickedData, log)
		}
	}

	fmt.Println("openedData length:", len(openedData))
	fmt.Println("clickedData length:", len(clickedData))

	openedDataresult := aggregateDataByDate(openedData, "openedAt")
	clickedDataresult := aggregateDataByDate(clickedData, "clickedAt")
	// Combine all unique dates from both opened and clicked data
	dateSet := make(map[string]struct{})
	for date := range openedDataresult {
		dateSet[date] = struct{}{}
	}
	for date := range clickedDataresult {
		dateSet[date] = struct{}{}
	}

	var allDates []string
	for date := range dateSet {
		allDates = append(allDates, date)
	}
	sort.Strings(allDates)

	openedSeriesData := make([]int, len(allDates))
	clickedSeriesData := make([]int, len(allDates))

	for i, date := range allDates {
		openedSeriesData[i] = openedDataresult[date]
		clickedSeriesData[i] = clickedDataresult[date]
	}

	type Response struct {
		AllDates          []string `json:"allDates"`
		OpenedSeriesData  []int    `json:"openedSeriesData"`
		ClickedSeriesData []int    `json:"clickedSeriesData"`
	}
	response := Response{
		AllDates:          allDates,
		OpenedSeriesData:  openedSeriesData,
		ClickedSeriesData: clickedSeriesData,
	}

	// Return the response
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}
func (db Database) updateScatterChartData(ctx *gin.Context) {

	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var logs []domains.TrackingLog

	// Fetch logs
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	//filter data
	var openedData []map[string]interface{}
	var clickedData []map[string]interface{}

	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			openedData = append(openedData, map[string]interface{}{
				"x":              log.OpenedAt.Unix() * 1000, // Convert to milliseconds
				"y":              log.ClickCount,
				"recipientEmail": log.RecipientEmail,
			})
		}
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			clickedData = append(clickedData, map[string]interface{}{
				"x":              log.ClickedAt.Unix() * 1000, // Convert to milliseconds
				"y":              log.ClickCount,
				"recipientEmail": log.RecipientEmail,
			})
		}
	}

	response := struct {
		OpenedData  []map[string]interface{} `json:"openedData"`
		ClickedData []map[string]interface{} `json:"clickedData"`
	}{
		OpenedData:  openedData,
		ClickedData: clickedData,
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)

}
func (db Database) barChartDataOpens(ctx *gin.Context) {

	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var logs []domains.TrackingLog
	//fetch logs
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	opensPerDay := make([]int, 7)
	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T01:00:00+01:00" {
			dayOfWeek := log.OpenedAt.Weekday() // Weekday() returns the day of the week (0 for Sunday, 1 for Monday, etc.)
			opensPerDay[dayOfWeek]++
		}
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, opensPerDay)

}
func (db Database) barChartDataClicks(ctx *gin.Context) {

	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}
	var logs []domains.TrackingLog
	//fetch logs
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}

	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	opensPerDay := make([]int, 7)
	for _, log := range logs {
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T01:00:00+01:00" {
			dayOfWeek := log.ClickedAt.Weekday() // Weekday() returns the day of the week (0 for Sunday, 1 for Monday, etc.)
			opensPerDay[dayOfWeek]++
		}
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, opensPerDay)
}
func (db Database) scatterChartDataOpens(ctx *gin.Context) {
	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Fetch logs
	var logs []domains.TrackingLog
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	opensPerHourDay := make(map[string]int)

	// Iterate through the logs and populate the map
	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T01:00:00+01:00" {
			date := log.OpenedAt
			hours := date.Hour()
			if date.Minute() >= 30 {
				hours++
			}
			dayOfWeek := int(date.Weekday())

			key := fmt.Sprintf("%d-%d", dayOfWeek, hours)
			opensPerHourDay[key]++
		}
	}
	type ScatterChartData struct {
		X int `json:"x"`
		Y int `json:"y"`
		Z int `json:"z"`
	}
	// Convert the map into a slice of structs for the scatter chart data
	var seriesData []ScatterChartData
	for key, value := range opensPerHourDay {
		var dayOfWeek, hour int
		fmt.Sscanf(key, "%d-%d", &dayOfWeek, &hour)
		seriesData = append(seriesData, ScatterChartData{
			X: dayOfWeek,
			Y: hour,
			Z: value,
		})
	}

	// Return the response
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, seriesData)

}
func (db Database) scatterChartDataClicks(ctx *gin.Context) {
	// Extract and validate companyID
	companyID := ctx.Param("companyID")
	if _, err := uuid.Parse(companyID); err != nil {
		logrus.Errorf("Invalid companyID format: %s, error: %v", companyID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid companyID format"})
		return
	}

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Fetch logs
	var logs []domains.TrackingLog
	if campaignID == "" {
		err := db.DB.Where("company_id = ?", companyID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	} else {
		err := db.DB.Where("company_id = ? AND campaign_id = ?", companyID, campaignID).Find(&logs).Error
		if err != nil {
			logrus.Error("Error occurred while finding logs. Error: ", err)
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return
		}
	}

	opensPerHourDay := make(map[string]int)

	// Iterate through the logs and populate the map
	for _, log := range logs {
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T01:00:00+01:00" {
			date := log.ClickedAt
			hours := date.Hour()
			if date.Minute() >= 30 {
				hours++
			}
			dayOfWeek := int(date.Weekday())

			key := fmt.Sprintf("%d-%d", dayOfWeek, hours)
			opensPerHourDay[key]++
		}
	}
	type ScatterChartData struct {
		X int `json:"x"`
		Y int `json:"y"`
		Z int `json:"z"`
	}
	// Convert the map into a slice of structs for the scatter chart data
	var seriesData []ScatterChartData
	for key, value := range opensPerHourDay {
		var dayOfWeek, hour int
		fmt.Sscanf(key, "%d-%d", &dayOfWeek, &hour)
		seriesData = append(seriesData, ScatterChartData{
			X: dayOfWeek,
			Y: hour,
			Z: value,
		})
	}

	// Return the response
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, seriesData)

}
