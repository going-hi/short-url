package link

import (
	"net/http"
	"short-url/pkg/jwt"
	"short-url/pkg/middleware"
)

type LinkHandlerParams struct {
	Repository *LinkRepository
	JwtService *jwt.JwtService
}

func NewLinkHandler(router *http.ServeMux, params LinkHandlerParams) {
	linkController := &LinkController{
		Repository: params.Repository,
	}

	router.Handle("POST /link", middleware.IsAuthMiddleware(http.HandlerFunc(linkController.create), params.JwtService))
	router.Handle("GET /link/{id}", middleware.IsAuthMiddleware(linkController.findById(), params.JwtService))
	router.Handle("GET /link", middleware.IsAuthMiddleware(linkController.getList(), params.JwtService))
	router.Handle("DELETE /link/{id}", middleware.IsAuthMiddleware(linkController.delete(), params.JwtService))
	router.Handle("GET /{code}", middleware.IsAuthMiddleware(linkController.GoTo(), params.JwtService))
}
