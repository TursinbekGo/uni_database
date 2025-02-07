package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create publication godoc
// @ID create_publication
// @Router /publication [POST]
// @Summary Create Publication
// @Description Create Publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param Publication body models.CreatePublication true "CreatePublicationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreatePublication(c *gin.Context) {

	var createPublication models.CreatePublication
	err := c.ShouldBindJSON(&createPublication)
	if err != nil {
		h.handlerResponse(c, "error Publication should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Publication().Create(c.Request.Context(), &createPublication)
	if err != nil {
		h.handlerResponse(c, "storage.Publication.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Publication().GetByID(c.Request.Context(), &models.PublicationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Publication.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Publication resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID publication godoc
// @ID get_by_id_publication
// @Router /publication/{id} [GET]
// @Summary Get By ID Publication
// @Description Get By ID Publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdPublication(c *gin.Context) {
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

	resp, err := h.strg.Publication().GetByID(c.Request.Context(), &models.PublicationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Publication.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Publication resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList publication godoc
// @ID get_list_publication
// @Router /publication [GET]
// @Summary Get List Publication
// @Description Get List Publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListPublication(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Publication offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Publication limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Publication().GetList(c.Request.Context(), &models.PublicationGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Publication.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Publication resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update publication godoc
// @ID update_publication
// @Router /publication/{id} [PUT]
// @Summary Update Publication
// @Description Update Publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Publication body models.UpdatePublication true "UpdatePublicationRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdatePublication(c *gin.Context) {

	var (
		id                string = c.Param("id")
		updatePublication models.UpdatePublication
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updatePublication)
	if err != nil {
		h.handlerResponse(c, "error Publication should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updatePublication.Id = id
	rowsAffected, err := h.strg.Publication().Update(c.Request.Context(), &updatePublication)
	if err != nil {
		h.handlerResponse(c, "storage.Publication.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Publication.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Publication().GetByID(c.Request.Context(), &models.PublicationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Publication.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Publication resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete publication godoc
// @ID delete_publication
// @Router /publication/{id} [DELETE]
// @Summary Delete Publication
// @Description Delete Publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeletePublication(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Publication().Delete(c.Request.Context(), &models.PublicationPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Publication.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Publication resposne", http.StatusNoContent, nil)
}

// GetPublicationStats godoc
// @ID get_publication_stats
// @Router /get_publication_stats [GET]
// @Summary Get  Publication Stats
// @Description Get the count of likes and downloads for a publication
// @Tags Publication
// @Accept json
// @Procedure json
// @Param publication_id query string true "Publication ID"
// @Success 200 {object} Response{data=models.PublicationStats} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetPublicationStats(c *gin.Context) {
	publicationID := c.Query("publication_id")

	fmt.Println("Received publication_id:", publicationID)
	if publicationID == "" {
		h.handlerResponse(c, "get list GetPublicationStats by publication_id", http.StatusBadRequest, "publication_id is required")
		return
	}
	if !helper.IsValidUUID(publicationID) {
		h.handlerResponse(c, "invalid publication ID", http.StatusBadRequest, "publication ID must be a valid UUID")
		return
	}

	stats, err := h.strg.Publication().GetPublicationStats(c.Request.Context(), publicationID)
	if err != nil {
		h.handlerResponse(c, "failed to get publication stats", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "publication stats retrieved successfully", http.StatusOK, stats)
}

// GetPublicationsByTag godoc
// @Router /publications/tags [GET]
// @Summary Get Publications by Tag
// @Description Get a list of publications that have the specified tag
// @Tags Publication
// @Accept json
// @Produce json
// @Param tag query string true "Tag to search for"
// @Success 200 {object} Response{data=[]models.Publication} "Success Request"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetPublicationsByTag(c *gin.Context) {
	tag := c.DefaultQuery("tag", "")
	if tag == "" {
		h.handlerResponse(c, "missing tag", http.StatusBadRequest, "Tag is required")
		return
	}

	publications, err := h.strg.Publication().GetPublicationsByTag(c.Request.Context(), tag)
	if err != nil {
		h.handlerResponse(c, "failed to get publications by tag", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "publications retrieved successfully", http.StatusOK, publications)
}
