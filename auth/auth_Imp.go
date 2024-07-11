package auth

import (
	"bioskuy/api/v1/user/entity"
	"bioskuy/helper"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	Env *helper.Config
}

func NewService(env *helper.Config) *jwtService {
	return &jwtService{
		Env: env,
	}
}

func (s *jwtService) GenerateToken(user entity.User, c *gin.Context) (string, error) {

	duration, _ := strconv.Atoi(s.Env.DurationJWT)

	claim := jwt.MapClaims{}
	claim["user_id"] = user.ID
	claim["name"] = user.Name
	claim["role"] = user.Role
	claim["email"] = user.Email
	claim["exp"] = time.Now().Add(time.Duration(duration) * time.Minute).Unix()


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString([]byte(s.Env.SecretKey))
	if err != nil {
		panic(err)
	}

	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(s.Env.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
