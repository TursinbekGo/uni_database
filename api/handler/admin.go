package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create admin godoc
// @ID create_admin
// @Router /admin [POST]
// @Summary Create Admin
// @Description Create Admin
// @Tags Admin
// @Accept json
// @Procedure json
// @Param Admin body models.CreateAdmin true "CreateAdminRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateAdmin(c *gin.Context) {

	var createAdmin models.CreateAdmin
	err := c.ShouldBindJSON(&createAdmin)
	if err != nil {
		h.handlerResponse(c, "error Admin should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Admin().Create(c.Request.Context(), &createAdmin)
	if err != nil {
		h.handlerResponse(c, "storage.Admin.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Admin().GetByID(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Admin.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Admin resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID admin godoc
// @ID get_by_id_admin
// @Router /admin/{id} [GET]
// @Summary Get By ID Admin
// @Description Get By ID Admin
// @Tags Admin
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdAdmin(c *gin.Context) {
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

	resp, err := h.strg.Admin().GetByID(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Admin.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Admin resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList admin godoc
// @ID get_list_admin
// @Router /admin [GET]
// @Summary Get List Admin
// @Description Get List Admin
// @Tags Admin
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListAdmin(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Admin offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Admin limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Admin().GetList(c.Request.Context(), &models.AdminGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Admin.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Admin resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update admin godoc
// @ID update_admin
// @Router /admin/{id} [PUT]
// @Summary Update Admin
// @Description Update Admin
// @Tags Admin
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Admin body models.UpdateAdmin true "UpdateAdminRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateAdmin(c *gin.Context) {

	var (
		id          string = c.Param("id")
		updateAdmin models.UpdateAdmin
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateAdmin)
	if err != nil {
		h.handlerResponse(c, "error Admin should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateAdmin.Id = id
	rowsAffected, err := h.strg.Admin().Update(c.Request.Context(), &updateAdmin)
	if err != nil {
		h.handlerResponse(c, "storage.Admin.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Admin.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Admin().GetByID(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Admin.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Admin resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete admin godoc
// @ID delete_admin
// @Router /admin/{id} [DELETE]
// @Summary Delete Admin
// @Description Delete Admin
// @Tags Admin
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteAdmin(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Admin().Delete(c.Request.Context(), &models.AdminPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Admin.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Admin resposne", http.StatusNoContent, nil)
}
