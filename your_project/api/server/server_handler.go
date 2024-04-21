package server

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

// CreateServer Handles the creation of a new server.
// @Summary        	Create server
// @Description    	Create a new server.
// @Tags			Servers
// @Accept			json
// @Produce			json
// @Param			request			body			ServerIn	true	"Server query params"
// @Success			201				{object}		utils.ApiResponses
// @Failure			400				{object}		utils.ApiResponses			"Invalid request"
// @Failure			401				{object}		utils.ApiResponses			"Unauthorized"
// @Failure			403				{object}		utils.ApiResponses			"Forbidden"
// @Failure			500				{object}		utils.ApiResponses			"Internal Server Error"
// @Router			:companyID/servers		[post]
func (db Database) CreateServer(ctx *gin.Context) {

	//Extract JWT values from the context
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
	// Parse the incoming JSON request into a ServerIn struct
	server := new(ServerIn)
	if err := ctx.ShouldBindJSON(server); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new server in the database
	dbServer := &domains.Server{
		ID:        uuid.New(),
		Name:      server.Name,
		CompanyID: companyID,
		Host:      server.Host,
		Port:      server.Port,
		Type:      server.Type,
		Username:  server.Username,
		Password:  server.Password,
	}

	if err := domains.Create(db.DB, dbServer); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, dbServer)
}

// ReadServers handles the retrieval of all servers.
// @Summary Get servers
// @Description Get all servers.
// @Tags Servers
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} ServerPaginator
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router :companyID/servers [get]
func (db Database) ReadServers(ctx *gin.Context) {
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

	// Retrieve all server data from the database
	servers, err := ReadAllPagination(db.DB, []domains.Server{}, companyID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all server data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count
	count, err := domains.ReadTotalCount(db.DB, &domains.Server{}, "company_id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a server structure as a response
	response := ServerPaginator{}
	listServer := []ServerTable{}
	for _, server := range servers {
		listServer = append(listServer, ServerTable{
			ID:        server.ID,
			Name:      server.Name,
			Host:      server.Host,
			Port:      server.Port,
			Type:      server.Type,
			Username:  server.Username,
			CreatedAt: server.CreatedAt,
		})
	}
	response.Items = listServer
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = count

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// ReadServer handles the retrieval of a specific server.
// @Summary Retrieve server
// @Description Get a specific server by ID.
// @Tags Servers
// @Produce json
// @Param ID path string true "Server ID"
// @Success 200 {object} domains.Server
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 404 {object} utils.ApiResponses "Not Found"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router :companyID/servers/:ID [get]
func (db Database) ReadServer(ctx *gin.Context) {
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
	// Parse and validate the server ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve the server data from the database
	var server domains.Server
	server, err = domains.ReadServerByID(db.DB, domains.Server{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving server data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Respond with the server data
	ctx.JSON(http.StatusOK, server)
}

// UpdateServer handles the update of a server.
// @Summary Update server
// @Description Update server.
// @Tags Servers
// @Accept json
// @Produce json
// @Param ID path string true "ServerID"
// @Param request body ServerIn true "Server query params"
// @Success 200 {object} utils.ApiResponses
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router :companyID/servers/:ID [put]
func (db Database) UpdateServer(ctx *gin.Context) {
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

	// Parse and validate the server ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a ServerIn struct
	server := new(ServerIn)
	if err := ctx.ShouldBindJSON(server); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the server with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Server{}, objectID); err != nil {
		logrus.Error("Error checking if the server with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the server data in the database
	dbServer := &domains.Server{
		Name:     server.Name,
		Host:     server.Host,
		Port:     server.Port,
		Type:     server.Type,
		Username: server.Username,
		Password: server.Password,
	}
	if err := domains.Update(db.DB, dbServer, objectID); err != nil {
		logrus.Error("Error updating server data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// DeleteServer handles the deletion of a server.
// @Summary Delete server
// @Description Delete one server.
// @Tags Servers
// @Produce json
// @Param ID path string true "Server ID"
// @Success 200 {object} utils.ApiResponses
// @Failure 400 {object} utils.ApiResponses "Invalid request"
// @Failure 401 {object} utils.ApiResponses "Unauthorized"
// @Failure 403 {object} utils.ApiResponses "Forbidden"
// @Failure 500 {object} utils.ApiResponses "Internal Server Error"
// @Router /servers/{ID} [delete]
func (db Database) DeleteServer(ctx *gin.Context) {
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
	// Parse and validate the server ID from the request parameter
	objectID, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Check if the server with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Server{}, objectID); err != nil {
		logrus.Error("Error checking if the server with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the server data from the database
	if err := domains.Delete(db.DB, &domains.Server{}, objectID); err != nil {
		logrus.Error("Error deleting server data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
