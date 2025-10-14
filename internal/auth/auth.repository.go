package auth

import (
	"database/sql"
	"short-url/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	Db *sql.DB
}

func (r *AuthRepository) FindByEmail(email string) (*user.User, error) {
	u := &user.User{}
	query := `SELECT id, email, password FROM users WHERE email=$1`
	err := r.Db.QueryRow(query, email).Scan(
		u.Id,
		u.Email,
		u.Password,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return u, nil
}

func (r *AuthRepository) Create(email, hashPassword string) (*user.User, error) {
	query := `INSERT INTO (email, password) VALUES ($1, $2) FROM users RETURNING id, email, password`

	u := &user.User{}

	err := r.Db.QueryRow(query, email, hashPassword).Scan(
		u.Id,
		u.Email,
		u.Password,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
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


