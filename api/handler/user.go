package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

//// @Security ApiKeyAuth
//// Create user godoc
//// @ID create_user
//// @Router /user [POST]
//// @Summary Create User
//// @Description Create User
//// @Tags User
//// @Accept json
//// @Procedure json
//// @Param User body models.CreateUser true "CreateUserRequest"
//// @Success 200 {object} Response{data=string} "Success Request"
//// @Response 400 {object} Response{data=string} "Bad Request"
//// @Failure 500 {object} Response{data=string} "Server error"

func (h *handler) CreateUser(c *gin.Context) {

	var createUser models.CreateUser
	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		h.handlerResponse(c, "error User should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.User().Create(c.Request.Context(), &createUser)
	if err != nil {
		h.handlerResponse(c, "storage.User.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.User().GetByID(c.Request.Context(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.User.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create User resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID user godoc
// @ID get_by_id_user
// @Router /user/{id} [GET]
// @Summary Get By ID User
// @Description Get By ID User
// @Tags User
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdUser(c *gin.Context) {
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

	resp, err := h.strg.User().GetByID(c.Request.Context(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.User.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id User resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList user godoc
// @ID get_list_user
// @Router /user [GET]
// @Summary Get List User
// @Description Get List User
// @Tags User
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListUser(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list User offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list User limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.User().GetList(c.Request.Context(), &models.UserGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.User.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list User resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update user godoc
// @ID update_user
// @Router /user/{id} [PUT]
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param User body models.UpdateUser true "UpdateUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateUser(c *gin.Context) {

	var (
		id         string = c.Param("id")
		updateUser models.UpdateUser
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateUser)
	if err != nil {
		h.handlerResponse(c, "error User should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateUser.Id = id
	rowsAffected, err := h.strg.User().Update(c.Request.Context(), &updateUser)
	if err != nil {
		h.handlerResponse(c, "storage.User.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.User.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.User().GetByID(c.Request.Context(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.User.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create User resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete user godoc
// @ID delete_user
// @Router /user/{id} [DELETE]
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteUser(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.User().Delete(c.Request.Context(), &models.UserPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.User.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create User resposne", http.StatusNoContent, nil)
}

// GetUserActivityCounts godoc
// @ID get_user_activity_counts
// @Router /get_user_activity_counts [GET]
// @Summary Get User Activity Counts
// @Description Get the count of publications, downloads, and likes for a specific user
// @Tags User
// @Accept json
// @Procedure json
// @Param user_id query string true "User ID"
// @Success 200 {object} Response{data=models.UserActivityCounts} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetUserActivityCounts(c *gin.Context) {

	userID := c.Query("user_id")

	fmt.Println("Received user_id:", userID)
	if userID == "" {
		h.handlerResponse(c, "get list GetUserActivityCounts by userID", http.StatusBadRequest, "user_id is required")
		return
	}
	if !helper.IsValidUUID(userID) {
		h.handlerResponse(c, "invalid user ID", http.StatusBadRequest, "user ID must be a valid UUID")
		return
	}

	counts, err := h.strg.User().GetUserActivityCounts(c.Request.Context(), userID)
	if err != nil {
		h.handlerResponse(c, "failed to get user activity counts", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "user activity counts retrieved successfully", http.StatusOK, counts)
}

// @Security ApiKeyAuth
// GetTopContributors godoc
// @ID get_top_contributors
// @Router /top_contributors [GET]
// @Summary Get Top Contributors
// @Description Get the top contributors (users) based on the count of their publications
// @Tags User
// @Accept json
// @Procedure json
// @Success 200 {object} Response{data=[]models.UserpublicationCounts} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetTopContributors(c *gin.Context) {
	// limit := c.DefaultQuery("limit", "10")

	// Convert limit to integer
	// limitInt, err := strconv.Atoi(limit)
	// if err != nil || limitInt <= 0 {
	// 	h.handlerResponse(c, "Invalid limit", http.StatusBadRequest, "limit should be a positive integer")
	// 	return
	// }

	// Fetch the top contributors
	contributors, err := h.strg.User().GetTopContributors(c.Request.Context())
	if err != nil {
		h.handlerResponse(c, "storage.User.get_top_contributors", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "top contributors response", http.StatusOK, contributors)
}

// @Security ApiKeyAuth
// GetUserScores godoc
// @ID get_user_scores
// @Router /users/scores [GET]
// @Summary Get User Scores
// @Description Get users ranked by score (likes + downloads) / publications
// @Tags User
// @Accept json
// @Procedure json
// @Param limit query int false "Number of top users to return (default 10)"
// @Success 200 {object} Response{data=[]models.UserScore} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetUserScores(c *gin.Context) {
	// limitStr := c.DefaultQuery("limit", "10")
	// limit, err := strconv.Atoi(limitStr)
	// if err != nil || limit < 1 {
	// 	h.handlerResponse(c, "invalid limit value", http.StatusBadRequest, "limit must be a positive integer")
	// 	return
	// }

	scores, err := h.strg.User().GetUserScores(c.Request.Context())
	if err != nil {
		h.handlerResponse(c, "failed to get user scores", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "user scores retrieved successfully", http.StatusOK, scores)
}

// @Security ApiKeyAuth
// GetUserRank godoc
// @ID get_user_rank
// @Router /users/{user_id}/rank [GET]
// @Summary Get User Rank
// @Description Get the rank of a specific user based on their score
// @Tags User
// @Accept json
// @Procedure json
// @Param user_id query string true "User ID"
// @Success 200 {object} Response{data=models.UserRank} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetUserRank(c *gin.Context) {
	userID := c.Query("user_id")
	fmt.Println("c.Param(\"user_id\"):", userID)
	rankInfo, err := h.strg.User().GetUserRank(c.Request.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") { // Check for "user not found" error
			h.handlerResponse(c, "user not found", http.StatusNotFound, err.Error())
			return
		}
		h.handlerResponse(c, "failed to get user rank", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "user rank retrieved successfully", http.StatusOK, rankInfo)
}

// @Security ApiKeyAuth
// GetUserStatistics godoc
// @ID get_user_statistics
// @Router /users/statistics [GET]
// @Summary Get User Statistics
// @Description Get count of publications, downloads, likes and rank for a specific user
// @Tags User
// @Accept json
// @Procedure json
// @Param user_id query string true "User ID"
// @Success 200 {object} Response{data=models.UserStatistics} "User Statistics"
// @Failure 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetUserStatistics(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "")
	if userID == "" {
		h.handlerResponse(c, "user_id is required", http.StatusBadRequest, "please provide a user_id")
		return
	}

	// Get the user statistics
	stats, err := h.strg.User().GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		h.handlerResponse(c, "failed to get user statistics", http.StatusInternalServerError, err.Error())
		return
	}

	// Return the statistics of the user
	h.handlerResponse(c, "user statistics retrieved successfully", http.StatusOK, stats)
}
