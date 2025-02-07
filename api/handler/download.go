package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create download godoc
// @ID create_download
// @Router /download [POST]
// @Summary Create Download
// @Description Create Download
// @Tags Download
// @Accept json
// @Procedure json
// @Param Download body models.CreateDownload true "CreateDownloadRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateDownload(c *gin.Context) {

	var createDownload models.CreateDownload
	err := c.ShouldBindJSON(&createDownload)
	if err != nil {
		h.handlerResponse(c, "error Download should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Download().Create(c.Request.Context(), &createDownload)
	if err != nil {
		h.handlerResponse(c, "storage.Download.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Download().GetByID(c.Request.Context(), &models.DownloadPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Download.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Download resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID download godoc
// @ID get_by_id_download
// @Router /download/{id} [GET]
// @Summary Get By ID Download
// @Description Get By ID Download
// @Tags Download
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdDownload(c *gin.Context) {
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

	resp, err := h.strg.Download().GetByID(c.Request.Context(), &models.DownloadPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Download.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Download resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList download godoc
// @ID get_list_download
// @Router /download [GET]
// @Summary Get List Download
// @Description Get List Download
// @Tags Download
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListDownload(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Download offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Download limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Download().GetList(c.Request.Context(), &models.DownloadGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Download.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Download resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update download godoc
// @ID update_download
// @Router /download/{id} [PUT]
// @Summary Update Download
// @Description Update Download
// @Tags Download
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Download body models.UpdateDownload true "UpdateDownloadRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateDownload(c *gin.Context) {

	var (
		id             string = c.Param("id")
		updateDownload models.UpdateDownload
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateDownload)
	if err != nil {
		h.handlerResponse(c, "error Download should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateDownload.Id = id
	rowsAffected, err := h.strg.Download().Update(c.Request.Context(), &updateDownload)
	if err != nil {
		h.handlerResponse(c, "storage.Download.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Download.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Download().GetByID(c.Request.Context(), &models.DownloadPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Download.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Download resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete download godoc
// @ID delete_download
// @Router /download/{id} [DELETE]
// @Summary Delete Download
// @Description Delete Download
// @Tags Download
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteDownload(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Download().Delete(c.Request.Context(), &models.DownloadPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Download.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Download resposne", http.StatusNoContent, nil)
}
