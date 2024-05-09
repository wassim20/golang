package workflow

import (
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// CreateWorkflow Handles the creation of a new Workflow.
// @Summary        	Create Workflow
// @Description    	Create a new Workflow.
// @Tags			workflows
// @Accept			json
// @Produce			json
// @Param			request			body			WorkflowIn	true	"Workflow query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/{companyID}/workflow		[post]
func (db Database) CreateWorkflow(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Parse the incoming JSON request into a WorkflowIn struct
	workflow := new(WorkflowIn)
	if err := ctx.ShouldBindJSON(workflow); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new workflow in the database
	dbWorkflow := &domains.Workflow{
		ID:            uuid.New(),
		Name:          workflow.Name,
		Status:        "pending",
		Trigger:       workflow.Trigger,
		MailinglistID: workflow.MailinglistID,
		CompanyID:     companyID,
	}

	if err := domains.Create(db.DB, dbWorkflow); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadWorkflows 	Handles the retrieval of all workflows.
// @Summary        	Get workflows
// @Description    	Get all workflows.
// @Tags			workflows
// @Produce			json
// @Param			page			query		int					false	"Page"
// @Param			limit			query		int					false	"Limit"
// @Success			200				{object}	WorkflowsPagination
// @Failure			400				{object}	utils.ApiResponses			"Invalid request"
// @Failure			401				{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses			"Forbidden"
// @Failure			500				{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/{companyID}/workflow		[get]
func (db Database) ReadWorkflows(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	log.Println(session, companyID)

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
	workflows, err := ReadAllPagination(db.DB, []domains.Workflow{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Workflow{}, "id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a workflow structure as a response
	response := WorkflowsPagination{}
	listWorkflow := []WorkflowTable{}
	for _, workflow := range workflows {

		listWorkflow = append(listWorkflow, WorkflowTable{
			ID:            workflow.ID,
			Name:          workflow.Name,
			Status:        workflow.Status,
			CurrentStep:   workflow.CurrentStep,
			Trigger:       workflow.Trigger,
			MailinglistID: workflow.MailinglistID,
			CompanyID:     workflow.CompanyID,
			CreatedAt:     workflow.CreatedAt,
		})
	}
	response.Items = listWorkflow
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadWorkflow 		Handles the retrieval of one Workflow.
// @Summary        	Get Workflow
// @Description    	Get one Workflow.
// @Tags			workflows
// @Produce			json
// @Param			ID   			path      	string		true		"Workflow ID"
// @Success			200				{object}	WorkflowDetails
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/workflow/{workflowID}	[get]
func (db Database) ReadWorkflow(ctx *gin.Context) {

	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the workflow ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("workflowID"))
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

	// Retrieve the workflow data by ID from the database
	workflow, err := ReadByID(db.DB, domains.Workflow{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a workflow structure as a response
	details := WorkflowDetails{
		ID:            workflow.ID,
		Name:          workflow.Name,
		Status:        workflow.Status,
		CurrentStep:   workflow.CurrentStep,
		Trigger:       workflow.Trigger,
		MailinglistID: workflow.MailinglistID,
		CompanyID:     workflow.CompanyID,
		Actions:       workflow.Actions,
		CreatedAt:     workflow.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateWorkflow 	Handles the update of a workflow.
// @Summary        	Update workflow
// @Description    	Update workflow.
// @Tags			workflows
// @Accept			json
// @Produce			json
// @Param			ID   			path      		string						true	"Workflow ID"
// @Param			request			body			WorkflowIn		true	"Workflow query params"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses				"Invalid request"
// @Failure			401				{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses				"Forbidden"
// @Failure			500				{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/{companyID}/workflow/{workflowID}	[put]
func (db Database) UpdateWorkflow(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the workflow ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("workflowID"))
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

	// Parse the incoming JSON request into a WorkflowIn struct
	workflow := new(WorkflowIn)
	if err := ctx.ShouldBindJSON(workflow); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the workflow with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Workflow{}, objectID); err != nil {
		logrus.Error("Error checking if the workflow with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the workflow data in the database
	dbWorkflow := &domains.Workflow{
		Name:          workflow.Name,
		Trigger:       workflow.Trigger,
		MailinglistID: workflow.MailinglistID,
		// Add other fields as necessary
	}
	if err := domains.Update(db.DB, dbWorkflow, objectID); err != nil {
		logrus.Error("Error updating workflow data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteWorkflow 	Handles the deletion of a workflow.
// @Summary        	Delete workflow
// @Description    	Delete one workflow.
// @Tags			workflows
// @Produce			json
// @Param			ID   			path      		string		true			"Workflow ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/workflow/{workflowID}	[delete]
func (db Database) DeleteWorkflow(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the workflow ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("workflowID"))
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

	// Check if the workflow with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Workflow{}, objectID); err != nil {
		logrus.Error("Error checking if the workflow with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the workflow data from the database
	if err := domains.Delete(db.DB, &domains.Workflow{}, objectID); err != nil {
		logrus.Error("Error deleting workflow data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// start a worflow
// @Summary        	Start workflow
// @Description    	Start a workflow.
// @Tags			workflows
// @Produce			json
// @Param			ID   			path      		string		true			"Workflow ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/workflow/{workflowID}/start	[post]
func (db Database) StartWorkflow(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse and validate the workflow ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("workflowID"))
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

	//get workflow
	workflow := domains.Workflow{}
	err = db.DB.Transaction(func(db *gorm.DB) error {
		// Check if the workflow with the specified ID exists
		if err := domains.CheckByID(db, &domains.Workflow{}, objectID); err != nil {
			logrus.Error("Error checking if the workflow with the specified ID exists. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
			return err
		}
		// retrieve the workflow

		err := db.First(&workflow, objectID).Error
		if err != nil {
			logrus.Error("Error retrieving workflow data from the database. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	//get contacts from mailinglist
	contacts := []domains.Contact{}
	err = db.DB.Transaction(func(db *gorm.DB) error {
		// Check if the mailinglist with the specified ID exists
		if err := domains.CheckByID(db, &domains.Mailinglist{}, workflow.MailinglistID); err != nil {
			logrus.Error("Error checking if the mailinglist with the specified ID exists. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
			return err
		}
		// retrieve the contacts
		err := db.Where("mailinglist_id = ?", workflow.MailinglistID).Find(&contacts).Error
		if err != nil {
			logrus.Error("Error retrieving contacts data from the database. Error: ", err.Error())
			utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	for _, contact := range contacts {

		// Start the workflow smilutaneously for all contacts
		go func(db *gorm.DB, workflow domains.Workflow, ctx *gin.Context, contact domains.Contact) {
			if err := Start(db, workflow, ctx, contact); err != nil {
				logrus.Error("Error starting workflow. Error: ", err.Error())
				utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
				return
			}
		}(db.DB, workflow, ctx, contact)

	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
