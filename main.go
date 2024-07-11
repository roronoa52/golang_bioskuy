package main

import (
	genreroute "bioskuy/api/v1/genre/route"
	genretomovieroute "bioskuy/api/v1/genretomovie/route"
	movieroute "bioskuy/api/v1/movies/route"
	paymentRoute "bioskuy/api/v1/payment/route"
	seatroute "bioskuy/api/v1/seat/route"
	seatbookingroute "bioskuy/api/v1/seatbooking/route"
	showtimeroute "bioskuy/api/v1/showtime/route"
	studioroute "bioskuy/api/v1/studio/route"
	"bioskuy/api/v1/user/route"
	"bioskuy/app"
	"bioskuy/exception"
	"bioskuy/helper"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {

	var c *gin.Context
	router := gin.Default()
	validate := validator.New()
	config := helper.NewConfig(c)
	db := app.GetConnection(config)
	defer db.Close()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, fmt.Sprint("Welcome to bioksuy"))
	})

	router.Use(exception.ErrorHandler)

	route.UserRoute(router, validate, db, config)
	genreroute.GenreRoute(router, validate, db, config)
	movieroute.MovieRoute(router, validate, db, config)
	genretomovieroute.GenreToMovieRoute(router, validate, db, config)
	studioroute.StudioRoute(router, validate, db, config)
	seatroute.SeatRoute(router, validate, db)
	showtimeroute.ShowtimeRoute(router, validate, db, config)
	seatbookingroute.SeatBookingRoute(router, validate, db, config)
	paymentRoute.PaymentRoute(router, validate, db, config)
	
	err := router.Run(":3000")
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return
	}
}
