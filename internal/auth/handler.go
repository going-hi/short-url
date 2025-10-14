package auth

import (
	"net/http"
	"short-url/internal/user"
	"short-url/pkg/jwt"
	"short-url/pkg/utils"
)

type AuthHandlerParams struct {
	*user.UserRepository
	*jwt.JwtService
}

type AuthController struct {
	Service *AuthService
	JwtService *jwt.JwtService
	UserRepository *user.UserRepository
}

func NewAuthHandler(router *http.ServeMux, params AuthHandlerParams) {
	service := &AuthService{}

	controller := &AuthController{
		Service: service,
		JwtService: params.JwtService,
		UserRepository: params.UserRepository,
	}

	router.HandleFunc("POST /auth/login", controller.login)
	router.HandleFunc("POST /auth/register", controller.register)

	// router.Handle("POST /auth/login", IsAuthMiddleware(http.HandlerFunc(controller.login), params.JwtService))
	// router.HandleFunc("GET /auth/logout", controller.logout)
}

func (controller *AuthController) login(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[LoginRequest](r.Body)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	userData, err := controller.UserRepository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	if userData == nil {
		utils.SendJson(w, 400, "Пользователя с таким email не существует")
		return
	}

	isMatch := controller.Service.CheckPassword(payload.Password, userData.Password)

	if !isMatch {
		utils.SendJson(w, 400, "Неправильный логин или пароль")
		return
	}

	jwtToken, err := controller.JwtService.GenerateJwt(userData)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken: jwtToken,
	}

	utils.SendJson(w, 200, response)
}

func (controller *AuthController) register(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[RegisterRequest](r.Body)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	userData, err := controller.UserRepository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	if userData != nil {
		utils.SendJson(w, 400, "Пользователь с такой почтой уже существует")
		return
	}

	hashPassword, err := controller.Service.HashPassword(payload.Password)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	user, err := controller.UserRepository.Create(payload.Email, hashPassword)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	token, err := controller.JwtService.GenerateJwt(user)

	if err != nil {
		utils.SendJson(w, 400, err.Error())
		return
	}

	response := RegisterResponse{
		AccessToken: token,
	}

	utils.SendJson(w, 200, response)
}