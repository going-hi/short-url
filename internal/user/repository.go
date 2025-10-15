package user

import (
	"database/sql"
)

type UserRepository struct {
	Db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	u := &User{}
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

func (r *UserRepository) Create(email, hashPassword string) (*User, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, password`

	u := &User{}

	err := r.Db.QueryRow(query, email, hashPassword).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}
