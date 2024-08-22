package handler

import (
	"net/http"
	"strings"

	constants "github.com/NidzamuddinMuzakki/movies-abishar/common/constant"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/response"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/logger"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	commonModel "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/response/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/validator"
	"github.com/NidzamuddinMuzakki/movies-abishar/model"
	"github.com/NidzamuddinMuzakki/movies-abishar/service"
	"github.com/gin-gonic/gin"
)

type IMovies interface {
	CreateMovies(c *gin.Context)
	UpdateMovies(c *gin.Context)
	GetListMovies(c *gin.Context)
	GetDetailMovies(c *gin.Context)
	GetDetailMoviesUsers(c *gin.Context)

	VoteMovies(c *gin.Context)
	UnVoteMovies(c *gin.Context)

	GetVoteMovies(c *gin.Context)
	GetMostVoteMovies(c *gin.Context)
	GetMostViewedMovies(c *gin.Context)
	GetMostViewedGenre(c *gin.Context)
}

type movies struct {
	common          common.IRegistry
	serviceRegistry service.IRegistry
}

func NewMovies(common common.IRegistry, serviceRegistry service.IRegistry) IMovies {
	return &movies{
		common:          common,
		serviceRegistry: serviceRegistry,
	}
}

func (h movies) GetVoteMovies(c *gin.Context) {
	ctx := c.Request.Context()
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)

	data, err := h.serviceRegistry.GetMoviesService().GetVoteMovies(ctx, result.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "List Vote User",
		Data:         data,
	})

}

func (h movies) GetMostViewedGenre(c *gin.Context) {
	ctx := c.Request.Context()
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	if strings.ToLower(result.Username) != "admin" {
		c.JSON(http.StatusForbidden, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusForbidden,
			Message:      "anda bukan admin",
		})
		return
	}
	data, err := h.serviceRegistry.GetMoviesService().GetMostViewedGenre(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Most Genre",
		Data:         data,
	})

}

func (h movies) GetMostVoteMovies(c *gin.Context) {
	ctx := c.Request.Context()
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	if strings.ToLower(result.Username) != "admin" {
		c.JSON(http.StatusForbidden, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusForbidden,
			Message:      "anda bukan admin",
		})
		return
	}
	data, err := h.serviceRegistry.GetMoviesService().GetMostVoteMovies(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Most Vote Movies",
		Data:         data,
	})

}

func (h movies) GetMostViewedMovies(c *gin.Context) {
	ctx := c.Request.Context()
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	if strings.ToLower(result.Username) != "admin" {
		c.JSON(http.StatusForbidden, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusForbidden,
			Message:      "anda bukan admin",
		})
		return
	}
	data, err := h.serviceRegistry.GetMoviesService().GetMostViewedMovies(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Most Movie",
		Data:         data,
	})

}
func (h movies) GetDetailMovies(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailMoviesModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
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

	data, err := h.serviceRegistry.GetMoviesService().GetMoviesById(ctx, int(payload.Id))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Detail MOvie",
		Data:         data,
	})

}

func (h movies) VoteMovies(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailMoviesModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
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
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)

	err := h.serviceRegistry.GetMoviesService().MoviesVote(ctx, int(payload.Id), result.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Vote Movie",
	})

}

func (h movies) UnVoteMovies(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailMoviesModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
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
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)

	err := h.serviceRegistry.GetMoviesService().UnVoteMovies(ctx, int(payload.Id), result.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "UnVote Movie",
	})

}
func (h movies) GetDetailMoviesUsers(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetDetailViewMoviesModel
	)

	if err := c.ShouldBindUri(&payload); err != nil {
		logger.Error(ctx, err.Error(), err)
		c.AbortWithStatusJSON(response.BadRequest(ctx).WithMessage(err.Error()).ToHTTPCodeAndMap())
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
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)

	data, err := h.serviceRegistry.GetMoviesService().GetMoviesByIdUserView(ctx, int(payload.Id), result.Username, int(payload.Duration))
	if err != nil {
		c.JSON(http.StatusBadRequest, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:         200,
		Success:      true,
		MessageTitle: commonModel.StatusSuccess,
		Message:      "Detail Movie",
		Data:         data,
	})

}
func (h movies) CreateMovies(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)
		payload model.RequestMoviesModel
	)
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	if strings.ToLower(result.Username) != "admin" {
		c.JSON(http.StatusForbidden, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusForbidden,
			Message:      "anda bukan admin",
		})
		return
	}
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

	err = h.serviceRegistry.GetMoviesService().CreateMovies(ctx, payload)
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
		Message:      "Movies Created",
	})
}

func (h movies) UpdateMovies(c *gin.Context) {
	// const logCtx = "delivery.http.tnc.CreateTnC"
	var (
		ctx = c.Request.Context()
		// span    = h.common.GetSentry().StartSpan(ctx, logCtx)
		payload model.RequestUpdateMoviesModel
	)
	tokenString := c.GetHeader(constants.Authorization)
	realToken := strings.Split(tokenString, "Bearer ")[1]
	result := util.ReadDataToken(realToken)
	if strings.ToLower(result.Username) != "admin" {
		c.JSON(http.StatusForbidden, response.Response{
			Success:      false,
			MessageTitle: "Failed",
			Code:         http.StatusForbidden,
			Message:      "anda bukan admin",
		})
		return
	}
	err := c.ShouldBindUri(&payload)
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
	err = c.ShouldBind(&payload)
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

	err = h.serviceRegistry.GetMoviesService().UpdateMovies(ctx, payload)
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
		Message:      "Movies Updated",
	})
}

func (h movies) GetListMovies(c *gin.Context) {
	var (
		ctx     = c.Request.Context()
		payload model.RequestGetListMoviesModel
	)

	if err := c.ShouldBind(&payload); err != nil {
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

	if payload.Limit == 0 {
		payload.Limit = 10
	}

	if payload.Offset == 0 {
		payload.Offset = 1
	}

	list, totalCount, err := h.serviceRegistry.GetMoviesService().GetMoviesList(ctx, payload)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, commonModel.Response{
			Status:  commonModel.StatusError,
			Message: http.StatusText(http.StatusUnprocessableEntity),
		})
		return
	}

	totalPages, previousPage, nextPage := util.Pagination(int64(totalCount), int64(payload.Limit), int64(payload.Offset))

	c.JSON(http.StatusOK, commonModel.Response{
		Status:       commonModel.StatusSuccess,
		Message:      http.StatusText(http.StatusOK),
		Data:         list,
		TotalRecords: totalCount,
		CurrentPage:  payload.Offset,
		NextPage:     uint(nextPage),
		PreviousPage: uint(previousPage),
		TotalPages:   uint(totalPages),
	})

}
