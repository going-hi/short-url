package auth

import (
	"net/http"
	"short-url/internal/user"
	"short-url/pkg/jwt"
	"short-url/pkg/utils"
)

type AuthController struct {
	Service        *AuthService
	JwtService     *jwt.JwtService
	UserRepository *user.UserRepository
}

type AuthControllerParams struct {
	Service        *AuthService
	UserRepository *user.UserRepository
	JwtService *jwt.JwtService
}

func NewAuthController(params *AuthControllerParams) *AuthController {
	return &AuthController{
		Service: params.Service,
		UserRepository: params.UserRepository,
		JwtService: params.JwtService,
	}
}

func (controller *AuthController) login(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[LoginRequest](r.Body)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := controller.UserRepository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if userData == nil {
		utils.SendJson(w, http.StatusBadRequest, "Пользователя с таким email не существует")
		return
	}

	isMatch := controller.Service.CheckPassword(payload.Password, userData.Password)

	if !isMatch {
		utils.SendJson(w, http.StatusBadRequest, "Неправильный логин или пароль")
		return
	}

	jwtToken, err := controller.JwtService.GenerateJwt(userData)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	response := LoginResponse{
		AccessToken: jwtToken,
	}

	utils.SendJson(w,  http.StatusOK, response)
}

func (controller *AuthController) register(w http.ResponseWriter, r *http.Request) {
	payload, err := utils.GetBody[RegisterRequest](r.Body)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := controller.UserRepository.FindByEmail(payload.Email)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	if userData != nil {
		utils.SendJson(w, http.StatusBadRequest, "Пользователь с такой почтой уже существует")
		return
	}

	hashPassword, err := controller.Service.HashPassword(payload.Password)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := controller.UserRepository.Create(payload.Email, hashPassword)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := controller.JwtService.GenerateJwt(user)

	if err != nil {
		utils.SendJson(w, http.StatusBadRequest, err.Error())
		return
	}

	response := RegisterResponse{
		AccessToken: token,
	}

	utils.SendJson(w, http.StatusOK, response)
}
