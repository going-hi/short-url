package main

import (
	"database/sql"
	"fmt"
	"short-url/config"
	_ "github.com/lib/pq"
)

func main() {
	config := config.LoadConfig()

	// поменять port=%s на %d
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Db.Host, config.Db.Port, config.Db.User, config.Db.Password, config.Db.DbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

// api
