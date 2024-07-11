package controller

import (
	"bioskuy/api/v1/studio/dto"
	"bioskuy/api/v1/studio/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type studioControllerImpl struct {
	studioService service.StudioService
}

func NewStudioController(studioService service.StudioService) StudioController {
	return &studioControllerImpl{studioService: studioService}
}

func (controller *studioControllerImpl) Create( c *gin.Context) {

	ctx := c.Request.Context()
	studio := dto.CreateStudioRequest{}

	err := c.ShouldBind(&studio)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := controller.studioService.Create(ctx, studio, c)
	if err != nil {
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data: result,
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *studioControllerImpl) FindById(c *gin.Context){

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("studioId")

	result, err := controller.studioService.FindByID(ctx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	} else {
		response.ResponseCode = http.StatusOK
		response.Data = result

		c.JSON(http.StatusOK, response)	
	}
}

func (controller *studioControllerImpl) FindAll(c *gin.Context) {
    ctx := c.Request.Context()

    result, err := controller.studioService.FindAll(ctx, c)
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

func (ctl *studioControllerImpl) Update(c *gin.Context) {
	ctx := c.Request.Context()
	studio := dto.UpdateStudioRequest{}

	id := c.Param("studioId")

	err := c.ShouldBind(&studio)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	studio.ID = id

	result, err := ctl.studioService.Update(ctx, studio, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusOK,
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

func (ctl *studioControllerImpl) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	response := web.FormatResponse{}
	id := c.Param("studioId")

	err := ctl.studioService.Delete(ctx, id, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return

	} else{
		response.ResponseCode = http.StatusOK
		response.Data = "OK"
	
		c.JSON(http.StatusOK, response)
	}
}