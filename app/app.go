package app

import (
	"database/sql"
	"net/http"
	"short-url/config"
	"short-url/internal/auth"
	"short-url/internal/link"
	"short-url/internal/user"
	"short-url/pkg/jwt"
)

type App struct {
	Server *http.Server
}

func NewApp(db *sql.DB, config *config.Config) *App {

	router := http.NewServeMux()

	jwtService := jwt.NewJwtService(config.SecretKey)


	userRepository := user.NewUserRepository(db)
	linkRepository := link.NewLinkRepository(db)
	

	auth.NewAuthHandler(router, auth.AuthHandlerParams{
		JwtService:     jwtService,
		UserRepository: userRepository,
	})

	link.NewLinkHandler(router, link.LinkHandlerParams{
		Repository: linkRepository,
		JwtService: jwtService,
	})

	app := &App{
		Server: &http.Server{
			Addr:    ":" + config.AppPort,
			Handler: router,
		},
	}

	return app
}

func (app *App) StartServer() error {
	err := app.Server.ListenAndServe()
	return err
}
