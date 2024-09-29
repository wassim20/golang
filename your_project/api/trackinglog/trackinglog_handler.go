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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var result struct {
		OpenedCount  int64 `json:"opened"`
		ClickedCount int64 `json:"clicked"`
		ErrorCount   int64 `json:"error"`
	}

	// Create the query
	query := db.DB.Model(&domains.TrackingLog{}).
		Select("COUNT(CASE WHEN opened_at IS NOT NULL AND opened_at != '0001-01-01T00:00:00Z' THEN 1 ELSE NULL END) AS opened_count, "+
			"COUNT(CASE WHEN clicked_at IS NOT NULL THEN 1 ELSE NULL END) AS clicked_count, "+
			"COUNT(CASE WHEN error IS NOT NULL AND error !='' THEN 1 ELSE NULL END) AS error_count").
		Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		// Parse the dates to ensure they are in the correct format
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Adjust the query to filter by date range
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// Execute the query
	err := query.Scan(&result).Error
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var result struct {
		OpenedCount  int64 `json:"opened_count"`
		ClickedCount int64 `json:"clicked_count"`
		ErrorCount   int64 `json:"error_count"`
	}

	// Create the base query
	query := db.DB.Model(&domains.TrackingLog{}).
		Select("COUNT(CASE WHEN opened_at IS NOT NULL AND opened_at != '0001-01-01T00:00:00Z' THEN 1 ELSE NULL END) AS opened_count, "+
			"COUNT(CASE WHEN clicked_at IS NOT NULL THEN 1 ELSE NULL END) AS clicked_count, "+
			"COUNT(CASE WHEN error IS NOT NULL AND error != '' THEN 1 ELSE NULL END) AS error_count").
		Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		// Parse the dates to ensure they are in the correct format
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Adjust the query to filter by date range
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// Execute the query
	err := query.Scan(&result).Error
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var result struct {
		TotalLogs  int64 `json:"total_logs"`
		OpenedLogs int64 `json:"opened_logs"`
	}

	// Create the base query
	query := db.DB.Model(&domains.TrackingLog{}).
		Select("COUNT(*) AS total_logs, "+
			"COUNT(CASE WHEN opened_at IS NOT NULL AND opened_at != '0001-01-01T00:00:00Z' THEN 1 ELSE NULL END) AS opened_logs").
		Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		// Parse the dates to ensure they are in the correct format
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Adjust the query to filter by date range
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// Execute the query
	err := query.Scan(&result).Error
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Calculate opened percentage
	openedPercentage := float64(0)
	if result.TotalLogs > 0 {
		openedPercentage = float64(result.OpenedLogs) / float64(result.TotalLogs) * 100
	}

	// Build and send the response
	type Response struct {
		OpenedPercentage float64 `json:"openedPercentage"`
		OpenedLogs       int64   `json:"openedLogs"`
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, Response{OpenedPercentage: openedPercentage, OpenedLogs: result.OpenedLogs})
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
		if timestamp != nil && !timestamp.IsZero() && timestamp.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var logs []domains.TrackingLog

	// Create the base query
	query := db.DB.Model(&domains.TrackingLog{}).
		Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		// Parse the dates to ensure they are in the correct format
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Adjust the query to filter by date range
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// Fetch logs based on the constructed query
	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Separate logs into openedData and clickedData
	var openedData []domains.TrackingLog
	var clickedData []domains.TrackingLog

	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			openedData = append(openedData, log)
		}
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T00:00:00Z" {
			clickedData = append(clickedData, log)
		}
	}

	fmt.Println("openedData length:", len(openedData))
	fmt.Println("clickedData length:", len(clickedData))

	openedDataresult := aggregateDataByDate(openedData, "openedAt")
	clickedDataresult := aggregateDataByDate(clickedData, "clickedAt")
	totalData := len(logs)

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
		Totals            int      `json:"totals"`
	}
	response := Response{
		AllDates:          allDates,
		OpenedSeriesData:  openedSeriesData,
		ClickedSeriesData: clickedSeriesData,
		Totals:            totalData,
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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var logs []domains.TrackingLog

	// Create the base query
	query := db.DB.Model(&domains.TrackingLog{}).
		Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		// Parse the dates to ensure they are in the correct format
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Adjust the query to filter by date range
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	}

	// Fetch logs based on the constructed query
	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Filter data
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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var logs []domains.TrackingLog
	query := db.DB.Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Use the "opened_at" column to filter logs
		query = query.Where("opened_at BETWEEN ? AND ?", start, end)
	}

	// Fetch logs
	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if len(logs) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"data": []int{}}) // Return empty data if no logs found
		return
	}

	opensPerDay := make([]int, 7)
	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T01:00:00Z" {
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

	// Extract and validate campaignID
	campaignID := ctx.Query("campaignID")
	if campaignID != "" {
		if _, err := uuid.Parse(campaignID); err != nil {
			logrus.Errorf("Invalid campaignID format: %s, error: %v", campaignID, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaignID format"})
			return
		}
	}

	// Extract start and end dates
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	var logs []domains.TrackingLog
	query := db.DB.Where("company_id = ?", companyID)

	// Add filters based on campaignID
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}

	// Add date filters if start and end dates are provided
	if startDate != "" && endDate != "" {
		start, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			logrus.Errorf("Invalid start date format: %s, error: %v", startDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}

		end, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			logrus.Errorf("Invalid end date format: %s, error: %v", endDate, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Use the "clicked_at" column to filter logs
		query = query.Where("clicked_at BETWEEN ? AND ?", start, end)
	}

	// Fetch logs
	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	if len(logs) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"data": []int{}}) // Return empty data if no logs found
		return
	}

	clicksPerDay := make([]int, 7)
	for _, log := range logs {
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T01:00:00Z" {
			dayOfWeek := log.ClickedAt.Weekday() // Weekday() returns the day of the week (0 for Sunday, 1 for Monday, etc.)
			clicksPerDay[dayOfWeek]++
		}
	}

	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, clicksPerDay)
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

	// Extract and validate date range
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")
	var startDate, endDate time.Time
	var dateRangeErr error

	if startDateStr != "" {
		startDate, dateRangeErr = time.Parse(time.RFC3339, startDateStr)
		if dateRangeErr != nil {
			logrus.Errorf("Invalid startDate format: %s, error: %v", startDateStr, dateRangeErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format"})
			return
		}
	}
	if endDateStr != "" {
		endDate, dateRangeErr = time.Parse(time.RFC3339, endDateStr)
		if dateRangeErr != nil {
			logrus.Errorf("Invalid endDate format: %s, error: %v", endDateStr, dateRangeErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format"})
			return
		}
	}

	// Fetch logs with optional date filtering
	var logs []domains.TrackingLog
	query := db.DB.Where("company_id = ?", companyID)
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}
	if !startDate.IsZero() {
		query = query.Where("opened_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("opened_at <= ?", endDate)
	}

	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	opensPerHourDay := make(map[string]int)

	// Iterate through the logs and populate the map
	for _, log := range logs {
		if log.OpenedAt != nil && log.OpenedAt.Format(time.RFC3339) != "0001-01-01T01:00:00Z" {
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

	var values []int
	for _, value := range opensPerHourDay {
		values = append(values, value)
	}

	// Check if values slice is empty before finding max/min
	if len(values) == 0 {
		logrus.Warn("No tracking logs found for the given date range")
		ctx.JSON(http.StatusOK, gin.H{"data": []ScatterChartData{}}) // or any appropriate response
		return
	}

	max := findMax(values)
	min := findMin(values)
	midThreshold := min + (max-min)/3
	largeThreshold := min + 2*(max-min)/3

	// Convert the map into a slice of structs for the scatter chart data with normalized Z values
	var seriesData []ScatterChartData
	for key, value := range opensPerHourDay {
		var dayOfWeek, hour int
		fmt.Sscanf(key, "%d-%d", &dayOfWeek, &hour)

		var size int
		if value <= midThreshold {
			size = 1 // small
		} else if value <= largeThreshold {
			size = 2 // medium
		} else {
			size = 3 // large
		}

		seriesData = append(seriesData, ScatterChartData{
			X: dayOfWeek,
			Y: hour,
			Z: size, // Normalized size
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

	// Extract and validate date range
	startDateStr := ctx.Query("startDate")
	endDateStr := ctx.Query("endDate")
	var startDate, endDate time.Time
	var dateRangeErr error

	if startDateStr != "" {
		startDate, dateRangeErr = time.Parse(time.RFC3339, startDateStr)
		if dateRangeErr != nil {
			logrus.Errorf("Invalid startDate format: %s, error: %v", startDateStr, dateRangeErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format"})
			return
		}
	}
	if endDateStr != "" {
		endDate, dateRangeErr = time.Parse(time.RFC3339, endDateStr)
		if dateRangeErr != nil {
			logrus.Errorf("Invalid endDate format: %s, error: %v", endDateStr, dateRangeErr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format"})
			return
		}
	}

	// Fetch logs with optional date filtering
	var logs []domains.TrackingLog
	query := db.DB.Where("company_id = ?", companyID)
	if campaignID != "" {
		query = query.Where("campaign_id = ?", campaignID)
	}
	if !startDate.IsZero() {
		query = query.Where("clicked_at >= ?", startDate)
	}
	if !endDate.IsZero() {
		query = query.Where("clicked_at <= ?", endDate)
	}

	err := query.Find(&logs).Error
	if err != nil {
		logrus.Error("Error occurred while finding logs. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	opensPerHourDay := make(map[string]int)

	// Iterate through the logs and populate the map
	for _, log := range logs {
		if log.ClickedAt != nil && log.ClickedAt.Format(time.RFC3339) != "0001-01-01T01:00:00Z" {
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

	var values []int
	for _, value := range opensPerHourDay {
		values = append(values, value)
	}

	// Check if values slice is empty before finding max/min
	if len(values) == 0 {
		logrus.Warn("No tracking logs found for the given date range")
		ctx.JSON(http.StatusOK, gin.H{"data": []ScatterChartData{}}) // or any appropriate response
		return
	}

	max := findMax(values)
	min := findMin(values)
	midThreshold := min + (max-min)/3
	largeThreshold := min + 2*(max-min)/3

	// Convert the map into a slice of structs for the scatter chart data with normalized Z values
	var seriesData []ScatterChartData
	for key, value := range opensPerHourDay {
		var dayOfWeek, hour int
		fmt.Sscanf(key, "%d-%d", &dayOfWeek, &hour)

		var size int
		if value <= midThreshold {
			size = 1 // small
		} else if value <= largeThreshold {
			size = 2 // medium
		} else {
			size = 3 // large
		}

		seriesData = append(seriesData, ScatterChartData{
			X: dayOfWeek,
			Y: hour,
			Z: size, // Normalized size
		})
	}

	// Return the response
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, seriesData)
}

func findMax(values []int) int {
	max := values[0]
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max
}

func findMin(values []int) int {
	min := values[0]
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min
}
