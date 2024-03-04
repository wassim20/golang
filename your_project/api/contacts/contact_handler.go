package contacts

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

// CreateContact 	Handles the creation of a new contact.
// @Summary        	Create contact
// @Description    	Create a new contact.
// @Tags			Contacts
// @Accept			json
// @Produce			json
// @Param			request			body			ContactIn	true	"Contact query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/:companyID/mailinglist/:mailinglistID/contacts		[post]
func (db Database) CreateContact(ctx *gin.Context) {

	//Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
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

	// Parse the incoming JSON request into a ContactIn struct
	contact := new(ContactIn)
	if err := ctx.ShouldBindJSON(contact); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	//fmt.Println(session.CompanyID)
	// Create a new contact in the database
	dbContact := &domains.Contact{
		ID:             uuid.New(),
		Email:          contact.Email,
		Firstname:      contact.Firstname,
		Lastname:       contact.Lastname,
		PhoneNumber:    contact.PhoneNumber,
		FullName:       contact.FullName,
		MailinglistsID: mailinglistID,

		// set current time as creation time
	}

	if err := domains.Create(db.DB, dbContact); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, dbContact)
}

// ReadContacts handles the retrieval of all contacts.
// @Summary Get contacts
// @Description Get all contacts.
// @Tags Contacts
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} ContactPaginator
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /:companyID/mailinglist/:mailinglistID/contacts [get]
func (db Database) ReadContacts(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
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

	log.Println(session)

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

	// Retrieve all contact data from the database

	contacts, err := ReadAllContactsForMailingList(db.DB, mailinglistID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all contact data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Contact{}, "id", mailinglistID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a contact structure as a response
	response := ContactPaginator{}
	listContact := []ContactTable{}
	for _, contact := range contacts {
		listContact = append(listContact, ContactTable{
			ID:          contact.ID,
			Email:       contact.Email,
			Firstname:   contact.Firstname,
			Lastname:    contact.Lastname,
			PhoneNumber: contact.PhoneNumber,
			FullName:    contact.FullName,
			CreatedAt:   contact.CreatedAt,

			// Ensure Mailinglist and Tags are properly set
		})
	}
	response.Items = listContact
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadContact handles the retrieval of one Contact.
// @Summary Get Contact
// @Description Get one Contact.
// @Tags Contacts
// @Produce json
// @Param ID path string true "Contact ID"
// @Success 200 {object} ContactDetails
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /:companyID/mailinglist/:mailinglistID/contacts/{ID} [get]
func (db Database) ReadContact(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
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

	// Parse and validate the contact ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the contact of a specific mailinglist data by ID  from the database
	contact, err := ReadContactByID(db.DB, domains.Contact{}, objectID, mailinglistID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a contact structure as a response
	details := ContactDetails{
		ID:          contact.ID,
		Email:       contact.Email,
		Firstname:   contact.Firstname,
		Lastname:    contact.Lastname,
		PhoneNumber: contact.PhoneNumber,
		FullName:    contact.FullName,
		CreatedAt:   contact.CreatedAt,
		Tags:        contact.Tags,
		// Ensure Mailinglist and Tags are properly set
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateContact	Handles the update of a contact.
// @Summary        	Update contact
// @Description    	Update contact.
// @Tags			Contacts
// @Accept			json
// @Produce			json
// @Param			ID   			path      		string						true	"ContactID"
// @Param			request			body			ContactIn		true	"Contact query params"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses				"Invalid request"
// @Failure			401				{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses				"Forbidden"
// @Failure			500				{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/:companyID/mailinglist/:mailinglistID/contacts/{ID}	[put]
func (db Database) UpdateContact(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
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

	// Parse and validate the contact ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a ContactIn struct
	contact := new(ContactIn)
	if err := ctx.ShouldBindJSON(contact); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the contact with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Contact{}, objectID); err != nil {
		logrus.Error("Error checking if the contact with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the contact data in the database
	dbContact := &domains.Contact{

		Email:          contact.Email,
		Firstname:      contact.Firstname,
		Lastname:       contact.Lastname,
		PhoneNumber:    contact.PhoneNumber,
		FullName:       contact.FullName,
		MailinglistsID: mailinglistID,

		// set current time as creation time
	}
	if err := domains.Update(db.DB, dbContact, objectID); err != nil {
		logrus.Error("Error updating contact data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteContact 	Handles the deletion of a contact.
// @Summary        	Delete contact
// @Description    	Delete one contact.
// @Tags			Contacts
// @Produce			json
// @Param			ID   			path      		string		true			"Contact ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/:companyID/mailinglist/:mailinglistID/contacts/{ID}	[delete]
func (db Database) DeleteContact(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
	// if err != nil {
	// 	logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
	// 	utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
	// 	return
	// }

	// Check if the employee belongs to the specified mailinglist
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the contact with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Contact{}, objectID); err != nil {
		logrus.Error("Error checking if the contact with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Delete the contact data from the database
	if err := domains.Delete(db.DB, &domains.Contact{}, objectID); err != nil {
		logrus.Error("Error deleting contact data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
