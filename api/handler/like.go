package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create like godoc
// @ID create_like
// @Router /like [POST]
// @Summary Create Like
// @Description Create Like
// @Tags Like
// @Accept json
// @Procedure json
// @Param Like body models.CreateLike true "CreateLikeRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateLike(c *gin.Context) {

	var createLike models.CreateLike
	err := c.ShouldBindJSON(&createLike)
	if err != nil {
		h.handlerResponse(c, "error Like should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Like().Create(c.Request.Context(), &createLike)
	if err != nil {
		h.handlerResponse(c, "storage.Like.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Like().GetByID(c.Request.Context(), &models.LikePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Like.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Like resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID like godoc
// @ID get_by_id_like
// @Router /like/{id} [GET]
// @Summary Get By ID Like
// @Description Get By ID Like
// @Tags Like
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdLike(c *gin.Context) {
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

	resp, err := h.strg.Like().GetByID(c.Request.Context(), &models.LikePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Like.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Like resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList like godoc
// @ID get_list_like
// @Router /like [GET]
// @Summary Get List Like
// @Description Get List Like
// @Tags Like
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListLike(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Like offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Like limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Like().GetList(c.Request.Context(), &models.LikeGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Like.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Like resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update like godoc
// @ID update_like
// @Router /like/{id} [PUT]
// @Summary Update Like
// @Description Update Like
// @Tags Like
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Like body models.UpdateLike true "UpdateLikeRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateLike(c *gin.Context) {

	var (
		id         string = c.Param("id")
		updateLike models.UpdateLike
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateLike)
	if err != nil {
		h.handlerResponse(c, "error Like should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateLike.Id = id
	rowsAffected, err := h.strg.Like().Update(c.Request.Context(), &updateLike)
	if err != nil {
		h.handlerResponse(c, "storage.Like.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Like.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Like().GetByID(c.Request.Context(), &models.LikePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Like.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Like resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete like godoc
// @ID delete_like
// @Router /like/{id} [DELETE]
// @Summary Delete Like
// @Description Delete Like
// @Tags Like
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteLike(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Like().Delete(c.Request.Context(), &models.LikePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Like.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Like resposne", http.StatusNoContent, nil)
}
