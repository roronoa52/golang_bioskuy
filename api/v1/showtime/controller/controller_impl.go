package controller

import (
	"bioskuy/api/v1/showtime/dto"
	"bioskuy/api/v1/showtime/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type showtimeControllerImpl struct {
	Service service.ShowtimeService
}

func NewMovieController(service service.ShowtimeService) ShowtimeController {
	return &showtimeControllerImpl{Service: service}
}

func (ctrl *showtimeControllerImpl) Create(c *gin.Context){
	ctx := c.Request.Context()
	showtime := dto.ShowtimeRequest{}

	err := c.ShouldBind(&showtime)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := ctrl.Service.Create(ctx, showtime, c)
	if err != nil {
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data: result,
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *showtimeControllerImpl) FindById(c *gin.Context){

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("showtimeId")

	result, err := controller.Service.FindByID(ctx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	} else {
		response.ResponseCode = http.StatusOK
		response.Data = result

		c.JSON(http.StatusOK, response)	
	}
}

func (controller *showtimeControllerImpl) FindAll(c *gin.Context) {
    ctx := c.Request.Context()

    result, err := controller.Service.FindAll(ctx, c)
    if err != nil {
        c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return
    }

    response := web.FormatResponse{
        ResponseCode: http.StatusOK,
        Data:    result,
    }

    c.JSON(http.StatusOK, response)
}

func (ctl *showtimeControllerImpl) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	response := web.FormatResponse{}
	id := c.Param("showtimeId")

	err := ctl.Service.Delete(ctx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return

	} else{
		response.ResponseCode = http.StatusOK
		response.Data = "OK"
	
		c.JSON(http.StatusOK, response)
	}
}
