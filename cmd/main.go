package main

import (
	"fmt"
	"net/http"
	"short-url/config"
	"short-url/internal/auth"
	"short-url/internal/link"
	"short-url/internal/user"
	"short-url/pkg/database"
	"short-url/pkg/jwt"

	_ "github.com/lib/pq"
)

// также можно вынести в одельный файл инизиализацию di и сделать компактнее майн также можно убрать дублирование кода
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

	fmt.Println("Successfully connected!")

	router := http.NewServeMux()

	jwtService := &jwt.JwtService{
		SecretKey: config.SecretKey,
	}

	userRepository := &user.UserRepository{
		Db: db,
	}

	linkRepository := &link.LinkRepository{
		Db: db,
	}

	auth.NewAuthHandler(router, auth.AuthHandlerParams{
		JwtService:     jwtService,
		UserRepository: userRepository,
	})

	link.NewLinkHandler(router, link.LinkHandlerParams{
		Repository: linkRepository,
		JwtService: jwtService,
	})

	server := http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router,
	}

	fmt.Println("Server is listening on port " + config.AppPort)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Server error:", err) // вместо fmt для сервака лучше использовать log
	}
}
