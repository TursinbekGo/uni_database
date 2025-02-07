package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create notification godoc
// @ID create_notification
// @Router /notification [POST]
// @Summary Create Notification
// @Description Create Notification
// @Tags Notification
// @Accept json
// @Procedure json
// @Param Notification body models.CreateNotification true "CreateNotificationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateNotification(c *gin.Context) {

	var createNotification models.CreateNotification
	err := c.ShouldBindJSON(&createNotification)
	if err != nil {
		h.handlerResponse(c, "error Notification should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Notification().Create(c.Request.Context(), &createNotification)
	if err != nil {
		h.handlerResponse(c, "storage.Notification.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Notification().GetByID(c.Request.Context(), &models.NotificationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Notification.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Notification resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID notification godoc
// @ID get_by_id_notification
// @Router /notification/{id} [GET]
// @Summary Get By ID Notification
// @Description Get By ID Notification
// @Tags Notification
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdNotification(c *gin.Context) {
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

	resp, err := h.strg.Notification().GetByID(c.Request.Context(), &models.NotificationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Notification.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Notification resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList notification godoc
// @ID get_list_notification
// @Router /notification [GET]
// @Summary Get List Notification
// @Description Get List Notification
// @Tags Notification
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListNotification(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Notification offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Notification limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Notification().GetList(c.Request.Context(), &models.NotificationGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Notification.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Notification resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update notification godoc
// @ID update_notification
// @Router /notification/{id} [PUT]
// @Summary Update Notification
// @Description Update Notification
// @Tags Notification
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Notification body models.UpdateNotification true "UpdateNotificationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateNotification(c *gin.Context) {

	var (
		id                 string = c.Param("id")
		updateNotification models.UpdateNotification
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateNotification)
	if err != nil {
		h.handlerResponse(c, "error Notification should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateNotification.Id = id
	rowsAffected, err := h.strg.Notification().Update(c.Request.Context(), &updateNotification)
	if err != nil {
		h.handlerResponse(c, "storage.Notification.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Notification.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Notification().GetByID(c.Request.Context(), &models.NotificationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Notification.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Notification resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete notification godoc
// @ID delete_notification
// @Router /notification/{id} [DELETE]
// @Summary Delete Notification
// @Description Delete Notification
// @Tags Notification
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteNotification(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Notification().Delete(c.Request.Context(), &models.NotificationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Notification.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Notification resposne", http.StatusNoContent, nil)
}
