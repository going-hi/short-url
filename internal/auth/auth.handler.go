package auth

import (
	"net/http"
	"short-url/config"
	"short-url/pkg/utils"
)

type AuthHandlerParams struct {
	*config.Config
	*AuthRepository
	*JwtService
}

type AuthController struct {
	Repository *AuthRepository
	JwtService *JwtService
}

func NewAuthHandler(router *http.ServeMux, params AuthHandlerParams) {

	controller := &AuthController{
		Repository: params.AuthRepository,
	}

	// router.Handle("POST /auth/login", IsAuthMiddleware(http.HandlerFunc(controller.login), params.JwtService))
	router.HandleFunc("POST /auth/login", controller.login)
	router.HandleFunc("POST /auth/register", controller.register)
	// router.HandleFunc("GET /auth/logout", controller.logout)
}

func (handler *AuthController) login(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[LoginRequest](r.Body)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	userData, err := handler.Repository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	if userData == nil {
		utils.SendJson(w, 400, "Пользователя с таким email не существует")
		return
	}

	isMatch := handler.Repository.CheckPassword(payload.Password, userData.Password)

	if !isMatch {
		utils.SendJson(w, 400, "Неправильный логин или пароль")
		return
	}

	jwtToken, err := handler.JwtService.GenerateJwt(userData)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken: jwtToken,
	}

	utils.SendJson(w, 200, response)
}

func (handler *AuthController) register(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[RegisterRequest](r.Body)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	userData, err := handler.Repository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	if userData != nil {
		utils.SendJson(w, 400, "Пользователь с такой почтой уже существует")
		return
	}

	hashPassword, err := handler.Repository.HashPassword(payload.Password)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	user, err := handler.Repository.Create(payload.Email, hashPassword)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	token, err := handler.JwtService.GenerateJwt(user)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	response := RegisterResponse{
		AccessToken: token,
	}

	utils.SendJson(w, 200, response)
}

func (handler *AuthController) logout(w http.ResponseWriter, r *http.Request) {

}

func (handler *AuthController) refresh(w http.ResponseWriter, r *http.Request) {

}
