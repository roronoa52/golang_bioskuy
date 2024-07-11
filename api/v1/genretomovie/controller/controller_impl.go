package controller

import (
	"bioskuy/api/v1/genretomovie/dto"
	"bioskuy/api/v1/genretomovie/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type genretomovieControllerImpl struct {
	Service service.GenreToMovieService
}

func NewGenreToMovieController(service service.GenreToMovieService) GenretomovieController {
	return &genretomovieControllerImpl{Service: service}
}

func (controller *genretomovieControllerImpl) Create( c *gin.Context) {

	ctx := c.Request.Context()
	genretomovie := dto.CreateGenreToMovieRequest{}

	err := c.ShouldBind(&genretomovie)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	result, err := controller.Service.Create(ctx, genretomovie, c)
	if err != nil {
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data: result,
	}

	c.JSON(http.StatusOK, response)
}

func (controller *genretomovieControllerImpl) FindById(c *gin.Context){

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("genretomovieId")

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

func (controller *genretomovieControllerImpl) FindAll(c *gin.Context) {
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

func (ctl *genretomovieControllerImpl) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	response := web.FormatResponse{}
	id := c.Param("genretomovieId")

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