package controller

import (
	"bioskuy/api/v1/seatbooking/dto"
	"bioskuy/api/v1/seatbooking/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type seatbookingControllerImpl struct {
	Service service.SeatBookingService
}

func NewSeatbookingController(service service.SeatBookingService) ShowtimeController {
	return &seatbookingControllerImpl{Service: service}
}

func (ctrl *seatbookingControllerImpl) Create(c *gin.Context) {
	ctx := c.Request.Context()
	showtime := dto.SeatBookingRequest{}

	err := c.ShouldBind(&showtime)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	userId := c.MustGet("user_id").(string)

	result, err := ctrl.Service.Create(ctx, showtime, userId, c)
	if err != nil {
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data:         result,
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *seatbookingControllerImpl) FindById(c *gin.Context) {

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("seatbookingId")

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

func (controller *seatbookingControllerImpl) FindAll(c *gin.Context) {
	ctx := c.Request.Context()

	result, err := controller.Service.FindAll(ctx, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusOK,
		Data:         result,
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *seatbookingControllerImpl) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	response := web.FormatResponse{}
	id := c.Param("seatbookingId")

	err := ctl.Service.Delete(ctx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return

	} else {
		response.ResponseCode = http.StatusOK
		response.Data = "OK"

		c.JSON(http.StatusOK, response)
	}
}
