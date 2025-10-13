package link

import (
	"database/sql"
)

type LinkRepository struct {
	db *sql.DB
}