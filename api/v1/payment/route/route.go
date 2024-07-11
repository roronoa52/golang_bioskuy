package movieroute

import (
	"bioskuy/api/v1/payment/controller"
	paymentRepo "bioskuy/api/v1/payment/repository"
	"bioskuy/api/v1/payment/service"
	seatRepo "bioskuy/api/v1/seat/repository"
	seatBookingRepo "bioskuy/api/v1/seatbooking/repository"
	"bioskuy/auth"
	"bioskuy/helper"
	"bioskuy/middleware"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func PaymentRoute(router *gin.Engine, validate *validator.Validate, db *sql.DB, config *helper.Config) {
	
	authService := auth.NewService(config)

	paymentRepo := paymentRepo.NewPaymentRepository()
	seatBookingRepo := seatBookingRepo.NewSeatBookingRepository()
	seatRepo := seatRepo.NewSeatRepository()
	paymentService := service.NewPaymentService(paymentRepo,seatRepo, seatBookingRepo, validate, db, config)
	paymentController := controller.NewPaymentController(paymentService)
	v1 := router.Group("/api/v1")
	{
		paymentRoutes := v1.Group("/payments")
		{
			paymentRoutes.POST("/", middleware.AuthMiddleware(authService, "user"), paymentController.Create)
			paymentRoutes.GET("/", paymentController.FindAll)
			paymentRoutes.GET("/:paymentId", paymentController.FindById)
			paymentRoutes.POST("/notification", paymentController.Notification)
		}
	}
}
