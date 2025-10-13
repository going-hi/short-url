package auth

type LoginRequest struct {
	Email    string `json:"name" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type RegisterRequest struct {
	Email    string `json:"name" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	AccessToken string `json:"accessToken"`
}
