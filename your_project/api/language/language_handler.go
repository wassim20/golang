package language

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

// @Summary      Create Language
// @Description  Creates a new language.
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Param        language body language.LanguageIn true "Language data"
// @Success      201       {object} utils.ApiResponses "Language created successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure      401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /languages [post]
func (db Database) CreateLanguage(ctx *gin.Context) {

	// Parse the incoming JSON request into a LanguageIn struct
	language := new(LanguageIn)
	if err := ctx.ShouldBindJSON(language); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Create a new language in the database
	dbLanguage := &domains.Language{
		ID:   uuid.New(),
		Name: language.Name,
		Code: language.Code,
	}

	if err := domains.Create(db.DB, dbLanguage); err != nil {
		logrus.Error("Error saving data to the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusCreated, constants.SUCCESS, utils.Null())
}

// @Summary      Get All Languages
// @Description  Retrieves a paginated list of languages.
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit"
// @Param        offset query int false "Offset"
// @Success      200       {object} utils.ApiResponses "Languages retrieved successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure	 401	   {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure 	500	   {object} utils.ApiResponses "Internal Server Error"
// @Router        /languages [get]
func (db Database) GetAllLanguages(ctx *gin.Context) {

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

	// Retrieve languages from the database
	languages, err := ReadAllPagination(db.DB, []domains.Language{}, limit, offset)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}
	// Retrieve total count of languages
	count, err := ReadTotalCount(db.DB)
	if err != nil {
		logrus.Error("Error occurred while finding total count. Error: ", err)
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Convert retrieved languages to response format
	response := LanguagePagination{}
	listLanguages := []LanguageTable{}

	for _, language := range languages {
		listLanguages = append(listLanguages, LanguageTable{
			Name: language.Name,
			Code: language.Code,
		})
	}
	response.Items = listLanguages
	response.TotalCount = uint(count)
	response.Page = uint(page)
	response.Limit = uint(limit)

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// @Summary      Get Language By ID
// @Description  Retrieves a language by its unique identifier.
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Param        id path string true "Language ID"
// @Success      200       {object} utils.ApiResponses "Language retrieved successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure	  	 401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Language not found"
// @Failure 	 500	   {object} utils.ApiResponses "Internal Server Error"
// @Router        /languages/:ID [get]
func (db Database) GetLanguageByID(ctx *gin.Context) {

	// Parse the language ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Retrieve language from the database
	language, err := ReadByID(db.DB, domains.Language{}, id)
	if err != nil {
		logrus.Error("Error retrieving data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Convert retrieved language to response format
	response := LanguageDetails{
		Name:      language.Name,
		Code:      language.Code,
		CreatedAt: language.CreatedAt,
		UpdatedAt: language.UpdatedAt,
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, response)
}

// @Summary      Update Language
// @Description  Updates an existing language by its unique identifier.
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Param        id path string true "Language ID"
// @Param        language body language.LanguageIn true "Language data"
// @Success      200       {object} utils.ApiResponses "Language updated successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure    401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Language not found"
// @Failure    500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /languages/:ID [put]
func (db Database) UpdateLanguage(ctx *gin.Context) {

	// Parse the language ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// Parse the incoming JSON request into a LanguageIn struct
	language := new(LanguageIn)
	if err := ctx.ShouldBindJSON(language); err != nil {
		logrus.Error("Error mapping request from frontend. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// check if the language exists
	if err := domains.CheckByID(db.DB, &domains.Language{}, id); err != nil {
		logrus.Error("Error checking if the Language with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Update the language details
	dbLanguage := &domains.Language{
		ID:   id,
		Name: language.Name,
		Code: language.Code,
	}
	if err := domains.Update(db.DB, dbLanguage, id); err != nil {
		logrus.Error("Error updating data in the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}

// @Summary      Delete Language
// @Description  Deletes an existing language by its unique identifier.
// @Tags         Languages
// @Accept       json
// @Produce      json
// @Param        id path string true "Language ID"
// @Success      200       {object} utils.ApiResponses "Language deleted successfully"
// @Failure      400       {object} utils.ApiResponses "Invalid request"
// @Failure    401       {object} utils.ApiResponses "Unauthorized"
// @Failure      403       {object} utils.ApiResponses "Forbidden"
// @Failure      404       {object} utils.ApiResponses "Language not found"
// @Failure    500       {object} utils.ApiResponses "Internal Server Error"
// @Router        /languages/:ID [delete]
func (db Database) DeleteLanguage(ctx *gin.Context) {

	// Parse the language ID from the request parameter
	id, err := uuid.Parse(ctx.Param("ID"))
	if err != nil {
		logrus.Error("Error mapping request from frontend. Invalid UUID format. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusBadRequest, constants.INVALID_REQUEST, utils.Null())
		return
	}

	// check if the language exists
	if err := domains.CheckByID(db.DB, &domains.Language{}, id); err != nil {
		logrus.Error("Error checking if the Language with the specified ID exists. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusNotFound, constants.DATA_NOT_FOUND, utils.Null())
		return
	}

	// Delete the language
	if err := domains.Delete(db.DB, &domains.Language{}, id); err != nil {
		logrus.Error("Error deleting data from the database. Error: ", err.Error())
		utils.BuildErrorResponse(ctx, http.StatusInternalServerError, constants.UNKNOWN_ERROR, utils.Null())
		return
	}

	// Respond with success
	utils.BuildResponse(ctx, http.StatusOK, constants.SUCCESS, utils.Null())
}
