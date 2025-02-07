package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create semester godoc
// @ID create_semester
// @Router /semester [POST]
// @Summary Create Semester
// @Description Create Semester
// @Tags Semester
// @Accept json
// @Procedure json
// @Param Semester body models.CreateSemester true "CreateSemesterRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateSemester(c *gin.Context) {

	var createSemester models.CreateSemester
	err := c.ShouldBindJSON(&createSemester)
	if err != nil {
		h.handlerResponse(c, "error Semester should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Semester().Create(c.Request.Context(), &createSemester)
	if err != nil {
		h.handlerResponse(c, "storage.Semester.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Semester().GetByID(c.Request.Context(), &models.SemesterPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Semester.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Semester resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID semester godoc
// @ID get_by_id_semester
// @Router /semester/{id} [GET]
// @Summary Get By ID Semester
// @Description Get By ID Semester
// @Tags Semester
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdSemester(c *gin.Context) {
	var id string = c.Param("id")

	// Here We Check id from Token
	// val, exist := c.Get("Auth")
	// if !exist {
	// 	h.handlerResponse(c, "Here", http.StatusInternalServerError, nil)
	// 	return
	// }

	// ToolData := val.(helper.TokenInfo)
	// if len(ToolData) > 0 {
	// 	id = ToolData.ToolID
	// } else {
	// 	id = c.Param("id")
	// }

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Semester().GetByID(c.Request.Context(), &models.SemesterPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Semester.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Semester resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList semester godoc
// @ID get_list_semester
// @Router /semester [GET]
// @Summary Get List Semester
// @Description Get List Semester
// @Tags Semester
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListSemester(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Semester offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Semester limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Semester().GetList(c.Request.Context(), &models.SemesterGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Semester.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Semester resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update semester godoc
// @ID update_semester
// @Router /semester/{id} [PUT]
// @Summary Update Semester
// @Description Update Semester
// @Tags Semester
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Semester body models.UpdateSemester true "UpdateSemesterRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateSemester(c *gin.Context) {

	var (
		id             string = c.Param("id")
		updateSemester models.UpdateSemester
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateSemester)
	if err != nil {
		h.handlerResponse(c, "error Semester should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateSemester.Id = id
	rowsAffected, err := h.strg.Semester().Update(c.Request.Context(), &updateSemester)
	if err != nil {
		h.handlerResponse(c, "storage.Semester.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Semester.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Semester().GetByID(c.Request.Context(), &models.SemesterPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Semester.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Semester resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete semester godoc
// @ID delete_semester
// @Router /semester/{id} [DELETE]
// @Summary Delete Semester
// @Description Delete Semester
// @Tags Semester
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteSemester(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Semester().Delete(c.Request.Context(), &models.SemesterPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Semester.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Semester resposne", http.StatusNoContent, nil)
}
