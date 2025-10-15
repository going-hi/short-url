package jwt

import (
	"short-url/internal/user"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func NewJwtService(secretKey string) *JwtService {
	return &JwtService{
		SecretKey: secretKey,
	}
}

type JwtService struct {
	SecretKey string
}

type JwtData struct {
	Id    int
	Email string
}

func (s *JwtService) GenerateJwt(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.Id,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.SecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *JwtService) VerifyJwt(tokenString string) (bool, *JwtData) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})

	if err != nil {
		return false, nil
	}

	if !token.Valid {
		return false, nil
	}

	email := token.Claims.(jwt.MapClaims)["email"]
	id := token.Claims.(jwt.MapClaims)["id"]

	return true, &JwtData{
		Id:    id.(int),
		Email: email.(string),
	}
}
