package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"app/api/models"
	"app/pkg/helper"
)

// Login godoc
// @ID login
// @Router /login [POST]
// @Summary Login
// @Description Login
// @Tags Login
// @Accept json
// @Procedure json
// @Param login body models.LoginInfo true "LoginRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) Login(c *gin.Context) {
	var login models.LoginInfo
	var role string
	var user_id string
	var passw string

	err := c.ShouldBindJSON(&login) // parse req body to given type struct
	if err != nil {
		h.handlerResponse(c, "create user", http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.strg.Admin().GetByID(context.Background(), &models.AdminPrimaryKey{Email: login.Email})
	if err == nil {
		user_id = admin.Id
		role = "admin"
		passw = admin.Password
		// hashedPassword = admin.Password
	} else {
		// Check in HeadNurse table
		user, err := h.strg.User().GetByID(context.Background(), &models.UserPrimaryKey{Email: login.Email})
		if err == nil {
			if user.Status == false {
				h.handlerResponse(c, "Account is inactive", http.StatusForbidden, "Account is inactive")
				return
			}
			user_id = user.Id
			role = "user"
			passw = user.Password
		} else {
			// No matching user found
			if err.Error() == "no rows in result set" {
				h.handlerResponse(c, "User does not exist", http.StatusBadRequest, "User does not exist")
				return
			}
			h.handlerResponse(c, "storage.user.getByID", http.StatusInternalServerError, err.Error())
			return
		}
	}

	if login.Password != passw {
		h.handlerResponse(c, "Wrong password", http.StatusBadRequest, "Wrong password")
		return
	}

	token, err := helper.GenerateJWT(map[string]interface{}{
		"user_id": user_id,
		"role":    role,
	}, time.Hour*24, h.cfg.SecretKey)

	h.handlerResponseLogin(c, "token", http.StatusCreated, token, role, user_id)

}

// Register godoc
// @ID register
// @Router /register [POST]
// @Summary Register
// @Description Register
// @Tags Register
// @Accept json
// @Procedure json
// @Param register body models.CreateUser true "CreateUserRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) Register(c *gin.Context) {

	var createUser models.CreateUser
	var id string
	err := c.ShouldBindJSON(&createUser)
	if err != nil {
		h.handlerResponse(c, "error user should bind json", http.StatusBadRequest, err.Error())
		return
	}

	if len(createUser.Password) < 7 {
		h.handlerResponse(c, "Password should inculude more than 7 elements", http.StatusBadRequest, errors.New("Password len should inculude more than 8 elements"))
		return
	}

	resp, err := h.strg.User().GetByID(context.Background(), &models.UserPrimaryKey{Email: createUser.Email})
	if err != nil {
		if err.Error() == "no rows in result set" {
			id, err = h.strg.User().Create(context.Background(), &createUser)
			if err != nil {
				h.handlerResponse(c, "storage.user.create", http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			h.handlerResponse(c, "User already exist", http.StatusInternalServerError, err.Error())
			return
		}
	} else if err == nil {
		h.handlerResponse(c, "User already exist", http.StatusBadRequest, nil)
		return
	}
	resp, err = h.strg.User().GetByID(context.Background(), &models.UserPrimaryKey{Id: id})

	h.handlerResponse(c, "create user resposne", http.StatusCreated, resp)
}
