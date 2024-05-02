package action

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateAction Handles the creation of a new Action.
// @Summary        	Create Action
// @Description    	Create a new Action.
// @Tags			actions
// @Accept			json
// @Produce			json
// @Param			request			body			ActionIn	true	"Action query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/:companyID/workflow/:workflowID/action		[post]
func (db Database) CreateAction(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	//ID for action to pass it to the send email in first case
	IDAction := uuid.New()
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	workflowID, err := uuid.Parse(ctx.Param("workflowID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Check if the employee belongs to the specified mailinglist
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a ActionIn struct
	action := new(ActionIn)
	if err := ctx.ShouldBindJSON(action); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Create a new action in the database
	dbAction := &domains.Action{
		ID:         IDAction,
		Name:       action.Name,
		ParentID:   action.ParentID,
		Type:       action.Type,
		Status:     "pending",
		WorkflowID: workflowID,
		Data:       action.Data,
	}

	if err := domains.Create(db.DB, dbAction); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadActions handles the retrieval of all actions.
// @Summary Get actions
// @Description Get all actions.
// @Tags Actions
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} ActionPaginator
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /:companyID/workflow/:workflowID/actions [get]
func (db Database) ReadActions(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	workflowID, err := uuid.Parse(ctx.Param("workflowID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified workflow
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

	// Retrieve all action data from the database
	actions, err := ReadAllPagination(db.DB, []domains.Action{}, workflowID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all action data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Action{}, "id", workflowID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate an action structure as a response
	response := ActionsPagination{}
	listAction := []ActionTable{}
	for _, action := range actions {
		listAction = append(listAction, ActionTable{
			ID:         action.ID,
			Name:       action.Name,
			Type:       action.Type,
			Status:     action.Status,
			WorkflowID: action.WorkflowID,
			Data:       action.Data,
			CreatedAt:  action.CreatedAt,
		})
	}
	response.Items = listAction
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadAction handles the retrieval of one Action.
// @Summary Get Action
// @Description Get one Action.
// @Tags actions
// @Produce json
// @Param ID path string true "Action ID"
// @Success 200 {object} ActionDetails
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /{companyID}/workflow/{workflowID}/action/{actionID} [get]
func (db Database) ReadAction(ctx *gin.Context) {

	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the workflow ID from the request parameter
	workflowID, err := uuid.Parse(ctx.Param("workflowID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the action ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("actionID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified workflow
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the action data by ID from the database
	action, err := ReadByID(db.DB, domains.Action{}, workflowID, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate an action structure as a response
	details := ActionDetails{
		ID:         action.ID,
		Name:       action.Name,
		Type:       action.Type,
		Status:     action.Status,
		WorkflowID: action.WorkflowID,
		Data:       action.Data,
		CreatedAt:  action.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateAction Handles the update of an action.
// @Summary Update action
// @Description Update action.
// @Tags actions
// @Accept json
// @Produce json
// @Param ID path string true "Action ID"
// @Param request body ActionIn true "Action query params"
// @Success 200 {object} utils.ApiResponses
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /{companyID}/workflow/{workflowID}/action/{actionID} [put]
func (db Database) UpdateAction(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the action ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("actionID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into an ActionIn struct
	action := new(ActionIn)
	if err := ctx.ShouldBindJSON(action); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the action with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Action{}, objectID); err != nil {
		logrus.Error("Error checking if the action with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the action data in the database
	dbAction := &domains.Action{
		Name: action.Name,
		Type: action.Type,
		Data: action.Data,
	}
	if err := domains.Update(db.DB, dbAction, objectID); err != nil {
		logrus.Error("Error updating action data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteAction Handles the deletion of an action.
// @Summary Delete action
// @Description Delete one action.
// @Tags actions
// @Produce json
// @Param ID path string true "Action ID"
// @Success 200 {object} utils.ApiResponses
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /{companyID}/workflow/{workflowID}/action/{actionID} [delete]
func (db Database) DeleteAction(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the action ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("actionID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified workflow
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the action with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Action{}, objectID); err != nil {
		logrus.Error("Error checking if the action with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the action data from the database
	if err := domains.Delete(db.DB, &domains.Action{}, objectID); err != nil {
		logrus.Error("Error deleting action data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
