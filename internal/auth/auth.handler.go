package auth

import (
	"encoding/json"
	"net/http"
	"short-url/config"
	"short-url/pkg/utils"
)

type AuthHandlerParams struct {
	*config.Config
	*AuthRepository
}

type AuthController struct {
	Repository *AuthRepository
}

func NewAuthHandler(router *http.ServeMux, params AuthHandlerParams) {

	controller := &AuthController{
		Repository: params.AuthRepository,
	}

	router.HandleFunc("POST /auth/login", controller.login)
	router.HandleFunc("POST /auth/register", controller.register)
	router.HandleFunc("GET /auth/logout", controller.logout)
}

func sendJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (handler *AuthController) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		sendJson(w, 400, err.Error())
		return
	}

	err = utils.IsValid(payload)

	if err != nil {
		sendJson(w, 400, err.Error())
		return
	}

	userData, err := handler.Repository.FindByEmail(payload.Email)

	if err != nil {
		sendJson(w, 400, err.Error())
		return
	}

	if userData == nil {
		sendJson(w, 400, "Пользователя с таким email не существует")
		return
	}

	isMatch := handler.Repository.CheckPassword(payload.Password, userData.Password)

	if !isMatch {
		sendJson(w, 400, "Неправильный логин или пароль")
		return
	}

	jwtToken, err := handler.Repository.GenerateJwt(userData)

	if err != nil {
		sendJson(w, 400, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken: jwtToken,
	}

	sendJson(w, 200, response)
}

func (handler *AuthController) register(w http.ResponseWriter, r *http.Request) {

}

func (handler *AuthController) logout(w http.ResponseWriter, r *http.Request) {

}

func (handler *AuthController) refresh(w http.ResponseWriter, r *http.Request) {

}
