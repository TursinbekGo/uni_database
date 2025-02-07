package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// @Security ApiKeyAuth
// Create course godoc
// @ID create_course
// @Router /course [POST]
// @Summary Create Course
// @Description Create Course
// @Tags Course
// @Accept json
// @Procedure json
// @Param Course body models.CreateCourse true "CreateCourseRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateCourse(c *gin.Context) {

	var createCourse models.CreateCourse
	err := c.ShouldBindJSON(&createCourse)
	if err != nil {
		h.handlerResponse(c, "error Course should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Course().Create(c.Request.Context(), &createCourse)
	if err != nil {
		h.handlerResponse(c, "storage.Course.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Course().GetByID(c.Request.Context(), &models.CoursePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Course.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Course resposne", http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// GetByID course godoc
// @ID get_by_id_course
// @Router /course/{id} [GET]
// @Summary Get By ID Course
// @Description Get By ID Course
// @Tags Course
// @Accept json
// @Procedure json
// @Param id path string false "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdCourse(c *gin.Context) {
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

	resp, err := h.strg.Course().GetByID(c.Request.Context(), &models.CoursePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Course.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Course resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList course godoc
// @ID get_list_course
// @Router /course [GET]
// @Summary Get List Course
// @Description Get List Course
// @Tags Course
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCourse(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Course offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Course limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Course().GetList(c.Request.Context(), &models.CourseGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Course.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Course resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// Update course godoc
// @ID update_course
// @Router /course/{id} [PUT]
// @Summary Update Course
// @Description Update Course
// @Tags Course
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Course body models.UpdateCourse true "UpdateCourseRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateCourse(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateCourse models.UpdateCourse
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateCourse)
	if err != nil {
		h.handlerResponse(c, "error Course should bind json", http.StatusBadRequest, err.Error())
		return
	}
	updateCourse.Id = id
	rowsAffected, err := h.strg.Course().Update(c.Request.Context(), &updateCourse)
	if err != nil {
		h.handlerResponse(c, "storage.Course.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Course.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Course().GetByID(c.Request.Context(), &models.CoursePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Course.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Course resposne", http.StatusAccepted, resp)
}

// @Security ApiKeyAuth
// Delete course godoc
// @ID delete_course
// @Router /course/{id} [DELETE]
// @Summary Delete Course
// @Description Delete Course
// @Tags Course
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteCourse(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Course().Delete(c.Request.Context(), &models.CoursePrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Course.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Course resposne", http.StatusNoContent, nil)
}

// GetList courses_by_semester_id godoc
// @ID get_list_courses_by_semester_id
// @Router /courses_by_semester_id [GET]
// @Summary Get List GetListCoursesBySemesterId
// @Description Get List GetListCoursesBySemesterId
// @Tags Course
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Param semester_id query string true "Semester ID" // New query parameter for semester_id
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCoursesBySemesterId(c *gin.Context) {
	semesterId := c.Query("semester_id")
	if semesterId == "" {
		h.handlerResponse(c, "get list PetrolHistory by semester_id", http.StatusBadRequest, "semester_id is required")
		return
	}
	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list PetrolHistory offset", http.StatusBadRequest, "invalid offset")
		return
	}

	// limit, err := h.getLimitQuery(c.Query("limit"))
	// if err != nil {
	// 	h.handlerResponse(c, "get list PetrolHistory limit", http.StatusBadRequest, "invalid limit")
	// 	return
	// }

	resp, err := h.strg.Course().GetList(c.Request.Context(), &models.CourseGetListRequest{
		Offset: offset,
		Limit:  100000,
		Search: c.Query("search"),
	})

	// cars, err := h.strg.Car().GetList(c.Request.Context(), &models.CarGetListRequest{
	// 	Offset: offset,
	// 	Limit:  limit,
	// 	Search: c.Query("search"),
	// })
	// if err != nil {
	// 	h.handlerResponse(c, "storage.Car.get_list", http.StatusInternalServerError, err.Error())
	// 	return
	// }
	var all []*models.Course

	for _, ptrh := range resp.Courses {
		if semesterId == ptrh.SemesterID {
			all = append(all, ptrh)
		}
	}

	if err != nil {
		h.handlerResponse(c, "storage.Semester.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Semester resposne", http.StatusOK, all)
}
