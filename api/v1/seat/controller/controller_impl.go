package controller

import (
	"bioskuy/api/v1/seat/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type seatControllerImpl struct {
	seatService service.SeatService
}

func NewSeatController(seatService service.SeatService) SeatController {
	return &seatControllerImpl{seatService: seatService}
}

func (controller *seatControllerImpl) FindById(c *gin.Context){

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("seatId")

	result, err := controller.seatService.FindByID(ctx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	} else {
		response.ResponseCode = http.StatusOK
		response.Data = result

		c.JSON(http.StatusOK, response)	
	}
}

func (controller *seatControllerImpl) FindAll(c *gin.Context) {
    ctx := c.Request.Context()
	id := c.Param("studioId")

    result, err := controller.seatService.FindAll(ctx, id, c)
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