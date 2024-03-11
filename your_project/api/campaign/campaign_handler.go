package campaign

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

// @Summary      Create Campaign
// @Description  Creates a new campaign.
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Param        campaign body campaign.CampaignIn true "Campaign data"
// @Success      201       {object} utils.ApiResponses "Campaign created successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure      401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      500       {object} utils.ApiResponses "Internal Server Error"
// @Router       :companyID/campaigns/:mailinglistID [post]
func (db Database) CreateCampaign(ctx *gin.Context) {

	// Extract JWT values from the context (assuming JWT middleware is used)
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
	fmt.Println("CompanyID", session.CompanyID, "pathID", companyID)
	// Check if the employee belongs to the specified company
	if err := domains.CheckEmployeeBelonging(db.DB, companyID, session.UserID, session.CompanyID); err != nil {
		logrus.Error("Error verifying employee belonging. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a CampaignIn struct
	campaign := new(CampaignIn)
	if err := ctx.ShouldBindJSON(campaign); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	fmt.Println("de", campaign.DeliveryAt)

	// Create a new campaign in the database
	dbCampaign := domains.Campaign{
		ID:              uuid.New(),
		CreatedByUserID: session.UserID,
		MailingListID:   mailinglistID,
		Type:            campaign.Type,
		Name:            campaign.Name,
		Subject:         campaign.Subject,
		HTML:            campaign.HTML,
		FromEmail:       campaign.FromEmail,
		FromName:        campaign.FromName,
		DeliveryAt:      campaign.DeliveryAt,
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic recovered:", r)
		}
	}()

	if err := domains.Create(db.DB, dbCampaign); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// @Summary      Get All Campaigns
// @Description  Retrieves a paginated list of all campaigns.
// @Tags         Campaigns
// @Produce      json
// @Param        page  query    int     false  "Page number" minimum(1)
// @Param        limit  query    int     false  "Results per page" minimum(1)
// @Success      200     {object}  campaign.CampaignsPagination "List of campaigns"
// @Failure      400     {object}  utils.ApiResponses          "Invalid request"
// @Failure      401     {object}  utils.ApiResponses          "Unauthorized"
// @Failure      403     {object}  utils.ApiResponses          "Forbidden"
// @Failure      500     {object}  utils.ApiResponses          "Internal Server Error"
// @Router       /:companyID/:mailinglistID/campaigns [get]
func (db Database) ReadAllCampaigns(ctx *gin.Context) {

	// Extract JWT values from the context (assuming JWT middleware is used)
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

	// Check if provided values are valid (minimum 1)
	if page < 1 || limit < 1 {
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Generate offset
	offset := (page - 1) * limit

	// Retrieve all campaign data from the database based on user's company
	campaigns, err := ReadAllPaginationFromMailinglist(db.DB, []domains.Campaign{}, mailinglistID, limit, offset)
	if err != nil {
		logrus.Error("Error occurred while finding all campaign data. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Retrieve total count of campaigns for the user's company
	count, err := domains.ReadTotalCount(db.DB, &domains.Campaign{}, "id", session.CompanyID)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Generate a campaign response structure
	response := CampaignsPagination{} // Consider renaming to CampaignsPagination
	listCampaigns := []CampaignsTable{}

	// Convert retrieved campaigns to response format
	for _, campaign := range campaigns {
		listCampaigns = append(listCampaigns, CampaignsTable{
			ID:          campaign.ID,
			Name:        campaign.Name,
			Subject:     campaign.Subject,
			FromEmail:   campaign.FromEmail,
			FromName:    campaign.FromName,
			ReplyTo:     campaign.ReplyTo,
			Status:      campaign.Status,
			SignDKIM:    campaign.SignDKIM,
			TrackOpen:   campaign.TrackOpen,
			TrackClick:  campaign.TrackClick,
			Resend:      campaign.Resend,
			CustomOrder: campaign.CustomOrder,
			RunAt:       campaign.RunAt,
			DeliveryAt:  campaign.DeliveryAt,
		})
	}
	response.Items = listCampaigns
	response.Page = uint(page)
	response.Limit = uint(limit)
	response.TotalCount = uint(count)

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// @Summary      Get Campaign
// @Description  Retrieves a specific campaign by its ID.
// @Tags         Campaigns
// @Produce      json
// @Param        id  path      string     true  "Campaign ID" format(uuid)
// @Success      200  {object}  campaign.CampaignsDetails  "Campaign details"
// @Failure      400  {object}  utils.ApiResponses          "Invalid request"
// @Failure      401  {object}  utils.ApiResponses          "Unauthorized"
// @Failure      403  {object}  utils.ApiResponses          "Forbidden"
// @Failure      404  {object}  utils.ApiResponses          "Campaign not found"
// @Failure      500  {object}  utils.ApiResponses          "Internal Server Error"
// @Router       /:companyID/campaigns/:ID [get]
func (db Database) ReadCampaign(ctx *gin.Context) {

	// Extract JWT values from the context (assuming JWT middleware is used)
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the company ID from the request parameter
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

	// Retrieve the campaign data by ID from the database
	campaign, err := ReadByID(db.DB, domains.Campaign{}, objectID)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	detailedCampaign := CampaignsDetails{
		ID:          campaign.ID,
		Name:        campaign.Name,
		Subject:     campaign.Subject,
		HTML:        campaign.HTML,
		Plain:       campaign.Plain,
		FromEmail:   campaign.FromEmail,
		FromName:    campaign.FromName,
		ReplyTo:     campaign.ReplyTo,
		Status:      campaign.Status,
		SignDKIM:    campaign.SignDKIM,
		TrackOpen:   campaign.TrackOpen,
		TrackClick:  campaign.TrackClick,
		Resend:      campaign.Resend,
		CustomOrder: campaign.CustomOrder,
		RunAt:       campaign.RunAt,
		DeliveryAt:  campaign.DeliveryAt,
	}
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, detailedCampaign)
}

// @Summary      Update Campaign
// @Description  Updates a specific campaign.
// @Tags         Campaigns
// @Accept       json
// @Produce      json
// @Param        id       path      string     true  "Campaign ID" format(uuid)
// @Param        campaign body      campaign.CampaignIn  true  "Campaign update data"
// @Success      200      {object}  utils.ApiResponses  "Campaign updated successfully"
// @Failure      400      {object}  utils.ApiResponses  "Invalid request"
// @Failure      401      {object}  utils.ApiResponses  "Unauthorized"
// @Failure      403      {object}  utils.ApiResponses  "Forbidden"
// @Failure      404      {object}  utils.ApiResponses  "Campaign not found"
// @Failure      500      {object}  utils.ApiResponses  "Internal Server Error"
// @Router       /:companyID/campaigns/:ID [put]
func (db Database) UpdateCampaign(ctx *gin.Context) {
	// Extract JWT values from the context (assuming JWT middleware is used)
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the campaign ID from the request parameter
	campaignID, err := uuid.Parse(ctx.Param("ID"))
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

	campaignUpdate := new(CampaignIn)
	if err := ctx.ShouldBindJSON(campaignUpdate); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}
	// Check if the campaign with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Campaign{}, campaignID); err != nil {
		logrus.Error("Error checking if the campaign with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	dbCampaign := domains.Campaign{
		ID:              campaignID,
		CreatedByUserID: session.UserID,
		Type:            campaignUpdate.Type,
		Name:            campaignUpdate.Name,
		Subject:         campaignUpdate.Subject,
		HTML:            campaignUpdate.HTML,
		FromEmail:       campaignUpdate.FromEmail,
		FromName:        campaignUpdate.FromName,
		DeliveryAt:      campaignUpdate.DeliveryAt,
	}
	if err := domains.Update(db.DB, dbCampaign, campaignID); err != nil {
		logrus.Error("Error updating company data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())

}

// @Summary      Delete Campaign
// @Description  Deletes a specific campaign.
// @Tags         Campaigns
// @Param        id       path      string     true  "Campaign ID" format(uuid)
// @Success      204      {object}  utils.ApiResponses  "Campaign deleted successfully"
// @Failure      400      {object}  utils.ApiResponses  "Invalid request"
// @Failure      401      {object}  utils.ApiResponses  "Unauthorized"
// @Failure      403      {object}  utils.ApiResponses  "Forbidden"
// @Failure      404      {object}  utils.ApiResponses  "Campaign not found"
// @Failure      500      {object}  utils.ApiResponses  "Internal Server Error"
// @Router       /:companyID/campaigns/:ID [delete]
func (db Database) DeleteCampaign(ctx *gin.Context) {

	// Extract JWT values from the context (assuming JWT middleware is used)
	session := utils.ExtractJWTValues(ctx)

	// Parse and validate the campaign ID from the request parameter
	campaignID, err := uuid.Parse(ctx.Param("ID"))
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

	// Check if the campaign with the specified ID exists
	if err := domains.CheckByID(db.DB, &domains.Campaign{}, campaignID); err != nil {
		logrus.Error("Error checking if the campaign with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the campaign data from the database
	if err := domains.Delete(db.DB, &domains.Campaign{}, campaignID); err != nil {
		logrus.Error("Error deleting company data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())

}
