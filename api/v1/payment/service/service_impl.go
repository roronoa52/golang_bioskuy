package service

import (
	"bioskuy/api/v1/payment/dto"
	"bioskuy/api/v1/payment/entity"
	"bioskuy/api/v1/payment/repository"
	entitySeat "bioskuy/api/v1/seat/entity"
	RepoSeat "bioskuy/api/v1/seat/repository"
	entitySeatBooking "bioskuy/api/v1/seatbooking/entity"
	RepoSeatBooking "bioskuy/api/v1/seatbooking/repository"
	"bioskuy/exception"
	"bioskuy/helper"
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type paymentServiceImpl struct {
	Repo repository.PaymentRepository
	RepoSeatBooking RepoSeatBooking.SeatBookingRepository
	RepoSeat RepoSeat.SeatRepository
	Validate *validator.Validate
	DB *sql.DB
	Env *helper.Config
}

func NewPaymentService(Repo repository.PaymentRepository, RepoSeat RepoSeat.SeatRepository, RepoSeatBooking RepoSeatBooking.SeatBookingRepository, validate *validator.Validate, DB *sql.DB, env *helper.Config) PaymentService {
	return &paymentServiceImpl{
		Repo: Repo,
		RepoSeatBooking: RepoSeatBooking,
		RepoSeat: RepoSeat,
		Validate: validate,
		DB: DB,
		Env: env,
	}
}

func (s *paymentServiceImpl)Create(ctx context.Context, request dto.PaymentRequest, userid string, c *gin.Context) (dto.CreatePaymentResponse, error)  {
	var PaymentResponse = dto.CreatePaymentResponse{}
	var PaymentRequest = dto.CreatePaymentRequest{}
	
	PaymentRequest.UserID = userid
	PaymentRequest.SeatDetailForBookingID = request.SeatDetailForBookingID

	err := s.Validate.Struct(request)
	if err != nil {
		c.Error(exception.ValidationError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return PaymentResponse, err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return PaymentResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	seatbookingExist, err := s.RepoSeatBooking.FindAllPendingByUserID(ctx, tx, PaymentRequest.UserID, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return PaymentResponse, err
	}

	totalSeat := len(seatbookingExist)
	total_price := totalSeat * seatbookingExist[0].MoviePrice

	payment := entity.Payment{
		UserID: userid,
		SeatDetailForBookingID: PaymentRequest.SeatDetailForBookingID,
		TotalSeat: totalSeat,
		TotalPrice: total_price,
	}

	result, err := s.Repo.Save(ctx, tx, payment, c)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return PaymentResponse, err
	}

	var servicemidtrans = snap.Client{}
	servicemidtrans.New(s.Env.MIDTRANS_SERVER_KEY, midtrans.Sandbox)

	req := & snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  result.ID,
			GrossAmt: int64(result.TotalPrice),
		}, 
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	snapResp, _ := servicemidtrans.CreateTransaction(req)

	fmt.Println("snap response: ", snapResp)

	PaymentResponse.ID = result.ID
	PaymentResponse.UserID = userid
	PaymentResponse.SeatDetailForBookingID = request.SeatDetailForBookingID
	PaymentResponse.TotalSeat = totalSeat
	PaymentResponse.TotalPrice = total_price
	PaymentResponse.URL = snapResp.RedirectURL

	return PaymentResponse, nil
}

func (s *paymentServiceImpl) FindByID(ctx context.Context, id string, c *gin.Context) (dto.PaymentResponse, error) {
	paymentResponse := dto.PaymentResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return paymentResponse, err
	}
	defer helper.CommitAndRollback(tx, c)

	result, err := s.Repo.FindByID(ctx, tx, id, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return paymentResponse, err
	}

	paymentResponse.ID = result.ID
	paymentResponse.SeatBookingStatus = result.SeatBookingStatus
	paymentResponse.UserID = result.UserID
	paymentResponse.ShowtimeID = result.ShowtimeID
	paymentResponse.ShowStart = result.ShowStart
	paymentResponse.ShowEnd = result.ShowEnd
	paymentResponse.StudioID = result.StudioID
	paymentResponse.StudioName = result.StudioName
	paymentResponse.MovieID = result.MovieID
	paymentResponse.MovieTitle = result.MovieTitle
	paymentResponse.MovieDescription = result.MovieDescription
	paymentResponse.MoviePrice = result.MoviePrice
	paymentResponse.MovieDuration = result.MovieDuration
	paymentResponse.MovieStatus = result.MovieStatus
	paymentResponse.SeatID = result.SeatID
	paymentResponse.SeatDetailForBookingID = result.SeatDetailForBookingID
	paymentResponse.SeatName = result.SeatName
	paymentResponse.SeatIsAvailable = result.SeatIsAvailable
	paymentResponse.TotalPrice = result.TotalPrice
	paymentResponse.TotalSeat = result.TotalSeat
	paymentResponse.Status = result.Status

	return paymentResponse, nil
}

func (s *paymentServiceImpl) FindAll(ctx context.Context, c *gin.Context) ([]dto.PaymentResponse, error) {
	paymentResponses := []dto.PaymentResponse{}

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return paymentResponses, err
	}
	defer helper.CommitAndRollback(tx, c)

	results, err := s.Repo.FindAll(ctx, tx, c)
	if err != nil {
		c.Error(exception.NotFoundError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return paymentResponses, err
	}

	for _, result := range results {
		paymentResponse := dto.PaymentResponse{
			ID:                    result.ID,
			SeatBookingStatus:     result.SeatBookingStatus,
			UserID:                result.UserID,
			ShowtimeID:            result.ShowtimeID,
			ShowStart:             result.ShowStart,
			ShowEnd:               result.ShowEnd,
			StudioID:              result.StudioID,
			StudioName:            result.StudioName,
			MovieID:               result.MovieID,
			MovieTitle:            result.MovieTitle,
			MovieDescription:      result.MovieDescription,
			MoviePrice:            result.MoviePrice,
			MovieDuration:         result.MovieDuration,
			MovieStatus:           result.MovieStatus,
			SeatID:                result.SeatID,
			SeatDetailForBookingID: result.SeatDetailForBookingID,
			SeatName:              result.SeatName,
			SeatIsAvailable:       result.SeatIsAvailable,
			TotalPrice: result.TotalPrice,
			TotalSeat: result.TotalSeat,
			Status: result.Status,
		}
		paymentResponses = append(paymentResponses, paymentResponse)
	}

	return paymentResponses, nil
}

func (s *paymentServiceImpl) Update(ctx context.Context, notificationPayload map[string]interface{}, c *gin.Context) {

	tx, err := s.DB.Begin()
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
	}
	defer helper.CommitAndRollback(tx, c)

	paymentId := notificationPayload["order_id"].(string)

	payment, _ := s.Repo.FindByID(ctx, tx, paymentId, c)

	fmt.Println(notificationPayload["order_id"])

	if notificationPayload["transaction_status"] == "settlement"{

		payment.Status = "paid"
		_, err := s.Repo.Update(ctx, tx, payment, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seatbooking := entitySeatBooking.SeatBooking{}
		seatbooking.SeatBookingStatus = "success"
		seatbooking.ID = payment.SeatBookingID

		_, err =  s.RepoSeatBooking.Update(ctx, tx, seatbooking, c )
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

	}else if notificationPayload["transaction_status"] == "deny"{

		err := s.Repo.Delete(ctx, tx, payment.ID, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seatbooking := entitySeatBooking.SeatBooking{}
		seatbooking.ID = payment.SeatBookingID
		
		err =  s.RepoSeatBooking.Delete(ctx, tx, seatbooking.ID, c )
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seat := entitySeat.Seat{}
		seat.ID = payment.SeatID
		seat.IsAvailable = true
		_, err = s.RepoSeat.Update(ctx, tx, seat, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

	}else if notificationPayload["transaction_status"] == "cancel" {

		err := s.Repo.Delete(ctx, tx, payment.ID, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seatbooking := entitySeatBooking.SeatBooking{}
		seatbooking.ID = payment.SeatBookingID
		
		err =  s.RepoSeatBooking.Delete(ctx, tx, seatbooking.ID, c )
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seat := entitySeat.Seat{}
		seat.ID = payment.SeatID
		seat.IsAvailable = true
		_, err = s.RepoSeat.Update(ctx, tx, seat, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

	}else if notificationPayload["transaction_status"] == "expire" {

		err := s.Repo.Delete(ctx, tx, payment.ID, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seatbooking := entitySeatBooking.SeatBooking{}
		seatbooking.ID = payment.SeatBookingID
		
		err =  s.RepoSeatBooking.Delete(ctx, tx, seatbooking.ID, c )
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}

		seat := entitySeat.Seat{}
		seat.ID = payment.SeatID
		seat.IsAvailable = true
		_, err = s.RepoSeat.Update(ctx, tx, seat, c)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		}
		
	}

}