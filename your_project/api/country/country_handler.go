package country

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

// @Summary      Create Country
// @Description  Creates a new country.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        country body country.CountryIn true "Country data"
// @Success      201       {object} utils.ApiResponses "Country created successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure      401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /countries [post]
func (db Database) CreateCountry(ctx *gin.Context) {

	// Parse the incoming JSON request into a CountryIn struct
	country := new(CountryIn)
	if err := ctx.ShouldBindJSON(country); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new country in the database
	dbCountry := &domains.Country{
		ID:        uuid.New(),
		Name:      country.Name,
		Code:      country.Code,
		Currency:  country.Currency,
		PhoneCode: country.PhoneCode,
		Flag:      country.Flag,
	}

	if err := domains.Create(db.DB, dbCountry); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// @Summary      Get All Countries
// @Description  Retrieves a paginated list of countries.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200       {object} utils.ApiResponses "Countries retrieved successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure	 401	   {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure 	500	   {object} utils.ApiResponses "Internal Server Error"
// @Router        /countries [get]
func (db Database) GetAllCountries(ctx *gin.Context) {

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

	// Retrieve countries from the database
	countries, err := ReadAllPagination(db.DB, []domains.Country{}, limit, offset)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Retrieve total count of campaigns for the user's company
	count, err := ReadTotalCount(db.DB)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Convert retrieved countries to response format
	response := CountryPagination{} // Consider renaming to CountryPagination
	listCountries := []CountryTable{}

	for _, country := range countries {
		listCountries = append(listCountries, CountryTable{
			Name:      country.Name,
			Code:      country.Code,
			Currency:  country.Currency,
			PhoneCode: country.PhoneCode,
			Flag:      country.Flag,
		})
	}
	response.Items = listCountries
	response.TotalCount = uint(count)
	response.Page = uint(page)
	response.Limit = uint(limit)

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// @Summary      Get Country By ID
// @Description  Retrieves a country by its unique identifier.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        id path string true "Country ID"
// @Success      200       {object} utils.ApiResponses "Country retrieved successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure	  	 401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Country not found"
// @Failure 	 500	   {object} utils.ApiResponses "Internal Server Error"
// @Router        /countries/:ID [get]
func (db Database) GetCountryByID(ctx *gin.Context) {

	// Parse the country ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve country from the database
	country, err := ReadByID(db.DB, domains.Country{}, id)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Convert retrieved country to response format
	response := CountryDetails{
		Name:      country.Name,
		Code:      country.Code,
		Currency:  country.Currency,
		PhoneCode: country.PhoneCode,
		Flag:      country.Flag,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// @Summary      Update Country
// @Description  Updates an existing country by its unique identifier.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        id path string true "Country ID"
// @Param        country body country.CountryIn true "Country data"
// @Success      200       {object} utils.ApiResponses "Country updated successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure    401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Country not found"
// @Failure    500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /countries/:ID [put]
func (db Database) UpdateCountry(ctx *gin.Context) {

	// Parse the country ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a CountryIn struct
	country := new(CountryIn)
	if err := ctx.ShouldBindJSON(country); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// check if the country exists
	if err := domains.CheckByID(db.DB, &domains.Country{}, id); err != nil {
		logrus.Error("Error checking if the Country with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the country details
	dbCountry := &domains.Country{
		ID:        id,
		Name:      country.Name,
		Code:      country.Code,
		Currency:  country.Currency,
		PhoneCode: country.PhoneCode,
		Flag:      country.Flag,
	}
	if err := domains.Update(db.DB, dbCountry, id); err != nil {
		logrus.Error("Error updating data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// @Summary      Delete Country
// @Description  Deletes an existing country by its unique identifier.
// @Tags         Countries
// @Accept       json
// @Produce      json
// @Param        id path string true "Country ID"
// @Success      200       {object} utils.ApiResponses "Country deleted successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure    401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Country not found"
// @Failure    500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /countries/:ID [delete]
func (db Database) DeleteCountry(ctx *gin.Context) {

	// Parse the country ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// check if the country exists
	if err := domains.CheckByID(db.DB, &domains.Country{}, id); err != nil {
		logrus.Error("Error checking if the Country with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the country
	if err := domains.Delete(db.DB, &domains.Country{}, id); err != nil {
		logrus.Error("Error deleting data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
