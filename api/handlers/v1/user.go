package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	models "testing-api-gateway/api/handlers/models"
	pbu "testing-api-gateway/genproto/user"
	l "testing-api-gateway/pkg/logger"
	"testing-api-gateway/pkg/utils"
)

// CreateUser ...
// @Summary CreateUser
// @Security ApiKeyAuth
// @Description Api for creating a new user
// @Tags user
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/createuser/ [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserSwagger
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockService().Create(ctx, &pbu.User{
		Id:       body.Id,
		Name:     body.Name,
		LastName: body.LastName,
		Username: body.Username,
		Email:    body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser gets user by id
// @Summary GetUser
// @Security ApiKeyAuth
// @Description Api for getting user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/getuser [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	id := c.Query("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockService().Get(
		ctx, &pbu.UserRequest{
			UserId: id,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListUsers returns list of users
// @Summary ListUsers
// @Security ApiKeyAuth
// @Description Api for getting users by page and limit
// @Tags user
// @Accept json
// @Produce json
// @Param Page path string true "page"
// @Param Limit path string true "limit"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [get]
func (h *handlerV1) ListUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockService().GetAll(
		ctx, &pbu.GetAllUsersRequest{
			Limit: params.Limit,
			Page:  params.Page,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser deletes user by id
// @Summary DeleteUser
// @Security ApiKeyAuth
// @Description Api for deleting users by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/deleteuser/ [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	var jsonMarshal protojson.MarshalOptions
	jsonMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.MockService().Delete(
		ctx, &pbu.UserRequest{
			UserId: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
