package main

import (
	"log"
	"short-url/app"
	"short-url/config"
	"short-url/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	config := config.LoadConfig()
	db, err := database.Connect(config.Db)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected!")


	app := app.NewApp(db, config)
	
	if err := app.StartServer(); err != nil {
		log.Fatalln("Server error:", err.Error())
	}

	log.Println("Server is listening on port " + config.AppPort)
}
