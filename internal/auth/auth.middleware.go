package auth

import (
	"context"
	"fmt"
	"net/http"
	"short-url/pkg/utils"
	"strings"
)

const (
	ContextEmailKey string = "ContextEmailKey"
	ContextIdKey    string = "ContextIdKey"
)

func IsAuthMiddleware(next http.Handler, service *JwtService) http.Handler {
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
		
		ctx := context.WithValue(r.Context(), ContextIdKey, jwtData.Id)
		ctx = context.WithValue(ctx, ContextEmailKey, jwtData.Email)
		req := r.WithContext(ctx)

		// next
		next.ServeHTTP(w, req) 

		fmt.Println("После запроса")
	})
}
