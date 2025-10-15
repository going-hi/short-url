package jwt

import (
	"short-url/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
	SecretKey string // используй лучше []byte для HMAC, jwt.SignedString ожидает байтовый ключ,сейчас не безопастно потому что ключ нельзя отчистить из памяти после использывания, таак как стринг имутабельный, под копотом у нас шифровка HMAC использует байты
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

func (s *JwtService) VerifyJwt(tokenString string) (bool, *JwtData) { // идимантически лучше возвращать (JwtData, error)
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
		Id:    id.(int), // используй float64 в int, прямое id.(int) может вызвать панику, числа из MapClaims всегда float64
		Email: email.(string),
	}
}
