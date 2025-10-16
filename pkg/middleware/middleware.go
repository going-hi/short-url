package middleware

import (
	"context"
	"fmt"
	"net/http"
	"short-url/pkg/jwt"
	"short-url/pkg/utils"
	"strings"
)

type ctxKey string

const (
	ContextEmailKey ctxKey = "ContextEmailKey"
	ContextIdKey    ctxKey = "ContextIdKey"
	ContextJWTData  ctxKey = "ContextJwtData"
)

func IsAuthMiddleware(next http.Handler, service *jwt.JwtService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			utils.SendJson(w, 401, "Неавторизован сука")
			return
		}

		authSplit := strings.Split(authorization, " ")

		// Улучшить
		typeToken := authSplit[0]
		token := authSplit[1]

		if typeToken == "" || token == "" {
			utils.SendJson(w, 401, "Неавторизован сука")
			return
		}

		if typeToken != "Bearer" {
			utils.SendJson(w, 401, "Неавторизован сука")
			return
		}

		isValid, jwtData := service.VerifyJwt(token)

		if !isValid {
			utils.SendJson(w, 401, "Не валидный токен")
			return
		}

		ctx := context.WithValue(r.Context(), ContextJWTData, jwtData)
		req := r.WithContext(ctx)

		// next
		next.ServeHTTP(w, req)

		fmt.Println("После запроса")
	})
}
