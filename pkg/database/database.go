package database

import (
	"fmt"
	"short-url/config"
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect(configDb config.Db) (*sql.DB, error) {
	// поменять port=%s на %d
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		configDb.Host, configDb.Port, configDb.User, configDb.Password, configDb.DbName)

	db, err := sql.Open("postgres", psqlInfo)

	return db, err
}