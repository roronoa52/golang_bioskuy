package controller

import (
	"bioskuy/api/v1/payment/dto"
	"bioskuy/api/v1/payment/service"
	"bioskuy/exception"
	"bioskuy/web"
	"net/http"

	"github.com/gin-gonic/gin"
)

type paymentControllerImpl struct {
	Service service.PaymentService
}

func NewPaymentController(service service.PaymentService) PaymentController {
	return &paymentControllerImpl{Service: service}
}

func (ctrl *paymentControllerImpl) Create(c *gin.Context){
	ctx := c.Request.Context()
	payment := dto.PaymentRequest{}

	err := c.ShouldBind(&payment)
	if err != nil {
		c.Error(exception.ForbiddenError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}

	userId := c.MustGet("user_id").(string)

	result, err := ctrl.Service.Create(ctx, payment, userId, c)
	if err != nil {
		return
	}

	response := web.FormatResponse{
		ResponseCode: http.StatusCreated,
		Data: result,
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *paymentControllerImpl) FindById(c *gin.Context){

	response := web.FormatResponse{}
	ctx := c.Request.Context()
	id := c.Param("paymentId")

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

func (controller *paymentControllerImpl) FindAll(c *gin.Context) {
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

func (controller *paymentControllerImpl) Notification(c *gin.Context){
    ctx := c.Request.Context()

	var notificationPayload map[string]interface{}

	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
        c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
        return
	}

	_, exists := notificationPayload["order_id"].(string)
	if !exists {
        c.Error(exception.ValidationError{Message: "order id not found"}).SetType(gin.ErrorTypePublic)
        return
	}

    controller.Service.Update(ctx, notificationPayload, c)

}
