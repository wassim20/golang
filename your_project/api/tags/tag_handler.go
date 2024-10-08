package tags

import (
	"fmt"
	"labs/constants"
	"labs/domains"
	"labs/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// CreateTag 	Handles the creation of a new Tag.
// @Summary        	Create Tag
// @Description    	Create a new Tag.
// @Tags			tags
// @Accept			json
// @Produce			json
// @Param			request			body			TagIn	true	"Tag query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			/{companyID}/Tags		[post]
func (db Database) CreateTag(ctx *gin.Context) {
	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	companyID, err := uuid.Parse(ctx.Param("companyID"))
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
	// Parse the incoming JSON request into a TagIn struct
	Tag := new(TagIn)
	if err := ctx.ShouldBindJSON(Tag); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, err.Error())
		return
	}

	// **New: Validate the TagIn struct using your function**
	if err := Validate_color(Tag); err != nil {
		logrus.Error("Error validating tag data: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, err.Error())
		return // Stop execution immediately
	}

	fmt.Println("Tag: ", Tag)
	// Create a new Tag in the database
	dbTag := &domains.Tag{
		ID:        uuid.New(),
		Name:      Tag.Name,
		Color:     Tag.Color,
		CompanyID: companyID,
	}

	fmt.Println("dbTag: ", dbTag)
	if err := domains.Create(db.DB, dbTag); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// ReadTagslist 	Handles the retrieval of all tags. with paginations
// @Summary        	Get tags
// @Description    	Get all tags.
// @Tags			tags
// @Produce			json
// @Param			page			query		int					false	"Page"
// @Param			limit			query		int					false	"Limit"
// @Success			200				{object}	TagPagination
// @Failure			400				{object}	utils.ApiResponses			"Invalid request"
// @Failure			401				{object}	utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses			"Forbidden"
// @Failure			500				{object}	utils.ApiResponses			"Internal Server Error"
// @Router			/{companyID}/tags		[get]
func (db Database) ReadTagslist(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Retrieve all tag data from the database
	tags, err := ReadAllPagination(db.DB, []domains.Tag{}, session.CompanyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all tag data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retriece total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Tag{}, "id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a tag structure as a response
	response := TagPagination{}
	listTags := []TagTable{}
	for _, tag := range tags {

		listTags = append(listTags, TagTable{
			ID:    tag.ID,
			Name:  tag.Name,
			Color: tag.Color,
		})
	}
	response.Items = listTags
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadTag 		Handles the retrieval of one Tag.
// @Summary        	Get Tag
// @Description    	Get one Tag.
// @Tags			tags
// @Produce			json
// @Param			ID   			path      	string		true		"Tag ID"
// @Success			200				{object}	TagDetails
// @Failure			400				{object}	utils.ApiResponses		"Invalid request"
// @Failure			401				{object}	utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}	utils.ApiResponses		"Forbidden"
// @Failure			500				{object}	utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/tags/{ID}	[get]
func (db Database) ReadTag(ctx *gin.Context) {

	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the tag ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("ID"))
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

	// Retrieve the tag data by ID from the database
	tag, err := ReadByID(db.DB, domains.Tag{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Generate a tag structure as a response
	details := TagDetails{
		ID:        tag.ID,
		Name:      tag.Name,
		Color:     tag.Color,
		CompanyID: tag.CompanyID,
		CreatedAt: tag.CreatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, details)
}

// UpdateTag 	Handles the update of a tag.
// @Summary        	Update tag
// @Description    	Update tag.
// @Tags			tags
// @Accept			json
// @Produce			json
// @Param			ID   			path      		string						true	"Tag ID"
// @Param			request			body			TagIn		true	"Tag query params"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses				"Invalid request"
// @Failure			401				{object}		utils.ApiResponses				"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses				"Forbidden"
// @Failure			500				{object}		utils.ApiResponses				"Internal Server Error"
// @Router			/{companyID}/tags/{ID}	[put]
func (db Database) UpdateTag(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the tag ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("companyID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the employee belongs to the specified Company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a TagIn struct
	tag := new(TagIn)
	if err := ctx.ShouldBindJSON(tag); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the tag data in the database
	dbTag := &domains.Tag{
		Name: tag.Name,
	}
	if err := domains.Update(db.DB, dbTag, objectID); err != nil {
		logrus.Error("Error updating tag data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, dbTag)
}

// DeleteTag 	Handles the deletion of a tag.
// @Summary        	Delete tag
// @Description    	Delete one tag.
// @Tags			tags
// @Produce			json
// @Param			ID   			path      		string		true			"Tag ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/tags/{ID}	[delete]
func (db Database) DeleteTag(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the tag ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the tag data from the database
	if err := domains.Delete(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error deleting tag data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// AssignTagToMailinglist 	Assign aTag To a Mailinglist
// @Summary        	Assaign tag to a mailinglist
// @Description    	Assaign one tag to a mailinglist.
// @Tags			tags
// @Produce			json
// @Param			ID   			path      		string		true			"Tag ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/tags/{ID}/mailinglist/{mailinglistID}	[POST]
func (db Database) AssignTagToMailinglist(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the tag ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Parse and validate the mailinglist ID from the request parameter
	mailinglistID, err := uuid.Parse(ctx.Param("mailinglistID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Check if the mailinglist with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Mailinglist{}, mailinglistID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// assign the tag data to the mailinglist
	if err := AssignToMailinglist(db.DB, objectID, mailinglistID); err != nil {
		logrus.Error("Error assigning a tag. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, "exists already in the mailinglist")
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// AssignTagToContact 	Assign aTag To a Contact
// @Summary        	Assaign tag to a contact
// @Description    	Assaign one tag to a contact.
// @Tags			tags
// @Produce			json
// @Param			ID   			path      		string		true			"Tag ID"
// @Success			200				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses		"Invalid request"
// @Failure			401				{object}		utils.ApiResponses		"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses		"Forbidden"
// @Failure			500				{object}		utils.ApiResponses		"Internal Server Error"
// @Router			/{companyID}/tags/{ID}/contact/{contactID}	[POST]
func (db Database) AssignTagToContact(ctx *gin.Context) {

	// Extract JWT values from the context
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the tag ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Parse and validate the mailinglist ID from the request parameter
	contactID, err := uuid.Parse(ctx.Param("contactID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Check if the mailinglist with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Contact{}, contactID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// assign the tag data to the mailinglist
	if err := AssignToContact(db.DB, objectID, contactID); err != nil {
		logrus.Error("Error assigning a tag. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

func (db Database) RemoveTagFromMailinglist(ctx *gin.Context) {

	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the mailinglist ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
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
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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

	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Check if the mailinglist with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Mailinglist{}, mailinglistID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Remove the tag from the mailinglist
	if err := RemoveFromMailinglist(db.DB, objectID, mailinglistID); err != nil {
		logrus.Error("Error removing tag from mailinglist. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, "error removing tag")
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

func (db Database) RemoveTagFromContact(ctx *gin.Context) {
	session := utils.ExtractJWTValues(ctx)
	// Parse and validate the mailinglist ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	contactID, err := uuid.Parse(ctx.Param("contactID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	companyID, err := uuid.Parse(ctx.Param("companyID"))
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
	// Check if the tag with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Tag{}, objectID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Check if the contact with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Contact{}, contactID); err != nil {
		logrus.Error("Error checking if the tag with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}
	// Remove the tag from the contact
	if err := RemoveFromContact(db.DB, objectID, contactID); err != nil {
		logrus.Error("Error removing tag from contact. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, "error removing tag")
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
