package auth

import (
	"net/http"
	"short-url/internal/user"
	"short-url/pkg/jwt"
)

type AuthHandlerParams struct {
	*user.UserRepository
	*jwt.JwtService
}

func NewAuthHandler(router *http.ServeMux, params AuthHandlerParams) {
	service := NewAuthService()

	controller := NewAuthController(
		&AuthControllerParams{
			Service:        service,
			JwtService:     params.JwtService,
			UserRepository: params.UserRepository,
		},
	)

	router.HandleFunc("POST /auth/login", controller.login)
	router.HandleFunc("POST /auth/register", controller.register)
}
