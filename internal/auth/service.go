package auth

import (
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// From docs bcrypt
func (r *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}


