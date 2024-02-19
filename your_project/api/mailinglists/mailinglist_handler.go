package mailinglists

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
)

// CreateMailinglist 	Handles the creation of a new Mailinglist.
// @Summary        	Create Mailinglist
// @Description    	Create a new Mailinglist.
// @Tags			Mailinglists
// @Accept			json
// @Produce			json

// @Param			request			body			mailinglists.MailinglistIn	true	"Mailinglist query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/mailinglist		[post]
func (db Database) CreateMailinglist(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse the incoming JSON request into a MailinglistIn struct
	mailinglist := new(MailinglistIn)
	if err := ctx.ShouldBindJSON(mailinglist); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new mailinglist in the database
	dbMailinglist := &domains.Mailinglist{
		ID:              uuid.New(),
		Name:            mailinglist.Name,
		Description:     mailinglist.Description,
		CompanyID:       session.CompanyID,
		CreatedByUserID: session.UserID,
		// set current time as creation time
	}

	if err := domains.Create(db.DB, dbMailinglist); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadMailinglists 	Handles the retrieval of all mailinglists.
// @Summary        	Get mailinglists
// @Description    	Get all mailinglists.
// @Tags			mailinglists
// @Produce			json

// @Param			page			query		int					false	"Page"
// @Param			limit			query		int					false	"Limit"
// @Success			200				{object}	MailinglistsPagination
// @Failure			400				{object}	utils.ApiResponses			"Invalid request"
// @Failure			401				{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses			"Forbidden"
// @Failure			500				{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/mailinglist		[get]
func (db Database) ReadMailinglists(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	log.Println("session")
	log.Println(session)
	log.Println("session")

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
	mailinglists, err := ReadAllPagination(db.DB, []domains.Mailinglist{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all company data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Mailinglist{}, "id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a mailinglist structure as a response
	response := MailinglistsPagination{}
	listMailinglist := []MailinglistTable{}
	for _, mailinglist := range mailinglists {

		listMailinglist = append(listMailinglist, MailinglistTable{
			ID:              mailinglist.ID,
			Name:            mailinglist.Name,
			Description:     mailinglist.Description,
			CompanyID:       mailinglist.CompanyID,
			CreatedByUserID: mailinglist.CreatedByUserID,
			CreatedAt:       mailinglist.CreatedAt,
		})
	}
	response.Items = listMailinglist
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadMailinglist 		Handles the retrieval of one Mailinglist.
// @Summary        	Get Mailinglist
// @Description    	Get one Mailinglist.
// @Tags			Companies
// @Produce			json

// @Param			ID   			path      	string		true		"Mailinglist ID"
// @Success			200				{object}	MailinglistDetails
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/mailinglist/{ID}	[get]
func (db Database) ReadMailinglist(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the mailinglist ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified mailinglist
	if err := domains.CheckEmployeeBelonging(db.DB, objectID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the mailinglist data by ID from the database
	mailinglist, err := ReadByID(db.DB, domains.Mailinglist{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a mailinglist structure as a response
	details := MailinglistDetails{
		ID:              mailinglist.ID,
		Name:            mailinglist.Name,
		Description:     mailinglist.Description,
		CompanyID:       mailinglist.CompanyID,
		CreatedByUserID: mailinglist.CreatedByUserID,
		CreatedAt:       mailinglist.CreatedAt,
		Contacts:        mailinglist.Contacts,
		Tags:            mailinglist.Tags,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}
