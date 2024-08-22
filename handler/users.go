package handler

import (
	"net/http"
	"strings"
	"time"

	constants "github.com/NidzamuddinMuzakki/movies-abishar/common/constant"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/response"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	commonCache "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/cache"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	commonModel "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/response/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/validator"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/service"
	"github.com/gin-gonic/gin"
)

type IUsers interface {
	CreateUsers(c *gin.Context)
	LoginUsers(c *gin.Context)
	LogoutUsers(c *gin.Context)
}

type users struct {
	common          common.IRegistry
	serviceRegistry service.IRegistry
}

func NewUsers(common common.IRegistry, serviceRegistry service.IRegistry) IUsers {
	return &users{
		common:          common,
		serviceRegistry: serviceRegistry,
	}
}

func (h users) CreateUsers(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)
		payload model.RequestUsersModel
	)

	err := c.ShouldBind(&payload)
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}
	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	err = h.serviceRegistry.GetUsersService().CreateUsers(ctx, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	tokens, err := util.GenerateTokenPair(payload.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusUnauthorized,
			Message:      err.Error(),
		})
		return
	}
	token := model.Token{
		Token: tokens["token"],
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Users Created",
		Data:         token,
	})
}

func (h users) LoginUsers(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)
		payload model.RequestUsersModel
	)

	err := c.ShouldBind(&payload)
	if err != nil {

		logger.Error(ctx, err.Error(), err)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}
	errValidate := h.common.GetValidator().Struct(payload)
	if errValidate != nil {
		dataError := validator.ToErrResponseV2(errValidate)
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      "Failed Data",
			Data:         dataError,
		})
		return
	}

	err = h.serviceRegistry.GetUsersService().LoginUsers(ctx, payload)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	tokens, err := util.GenerateTokenPair(payload.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusUnauthorized,
			Message:      err.Error(),
		})
		return
	}
	token := model.Token{
		Token: tokens["token"],
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Users Login",
		Data:         token,
	})
}

func (h users) LogoutUsers(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)

	)

	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	dataToken := commonCache.Data{
		Value: realToken,
		Key:   commonCache.Key(result.Uuid),
	}

	err := h.common.GetCache().Set(ctx, dataToken, time.Hour*2)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusBadRequest,
			Message:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Users Logout",
	})
}
