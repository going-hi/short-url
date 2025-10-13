package auth

import (
	"database/sql"
	"short-url/config"
	"short-url/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	db *sql.DB
	*config.Config
}

func (r *AuthRepository) FindByEmail(email string) (*user.User, error) {
	user := user.User{}
	query := `SELECT * FROM users WHERE email=$1`
	err := r.db.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}


func (r *AuthRepository) Create() {

}

// From docs bcrypt
func (r *AuthRepository) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *AuthRepository) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func (r *AuthRepository) GenerateJwt(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"id":    user.Id,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})


	tokenString, err := token.SignedString(r.SecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
