package trackinglog

import (
	"labs/api/campaign"
	"labs/constants"
	"labs/domains"
	"net/http"
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
// @Router       /{companyID}/{camapignID}/logs [post]
func (db Database) CreateTrackingLog(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("camapignID"))
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
// @Router       /{companyID}/{camapignID}/logs [get]
func (db Database) ReadTrackingLogs(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("camapignID"))
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
	logs, err := ReadAllPagination(db.DB, []domains.TrackingLog{}, session.CompanyID, campaignID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := ReadTotalCountTrackingLog(db.DB, session.CompanyID, campaignID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a log structure as a response
	response := TrackingLogPagination{}
	listlogs := []TrackingLogTable{}
	for _, log := range logs {

		listlogs = append(listlogs, TrackingLogTable{
			ID:         log.ID,
			CompanyID:  log.CompanyID,
			CampaignID: log.CampaignID,
			Status:     log.Status,
			Error:      log.Error,
			CreatedAt:  log.CreatedAt,
			UpdatedAt:  log.UpdatedAt,
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
// @Router			 /{companyID}/{camapignID}/logs/{ID}	[get]
func (db Database) ReadTrackingLogByID(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("camapignID"))
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
// @Router			/{companyID}/{camapignID}/logs/{ID}	[put]
func (db Database) UpdateTrackingLog(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	campaignID, err := uuid.Parse(ctx.Param("camapignID"))
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
// @Router       /{companyID}/{camapignID}/logs/{ID} [delete]
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
	// camapignID, err := uuid.Parse(ctx.Param("camapignID"))
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
// @Router       /{companyID}/{camapignID}/logs/open/{trackingID} [delete]
func (db Database) handleOpenRequest(ctx *gin.Context) {
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
	trackingLog.Status = "opened"     // Update status if needed
	trackingLog.OpenedAt = time.Now() // Update opened_at timestamp if needed
	trackingLog.ClickCount++          // Increment click count if needed

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
// @Router       /{companyID}/{camapignID}/logs/click/{trackingID} [delete]
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
	trackingLog.Status = "opened"     // Update status if needed
	trackingLog.OpenedAt = time.Now() // Update opened_at timestamp if needed
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
