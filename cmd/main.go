package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"short-url/config"
	"short-url/internal/auth"
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

	authRepository := &auth.AuthRepository{Db: db}

	router := http.NewServeMux()

	auth.NewAuthHandler(router, auth.AuthHandlerParams{
		Config: config,
		AuthRepository: authRepository,
	})
}

// api
