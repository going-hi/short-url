package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// var validate = validator.New() создаём валидатор один раз, не каждый вызов, так безопаснее и быстрее

func IsValid(payload any) error {
	var validate = validator.New() // вот тут убрать, и использовать уже созданный валидатор
	err := validate.Struct(payload)
	return err
}

func SendJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func GetBody[T any](body io.ReadCloser) (*T, error) {

	//defer body.Close() чтобы небыло утечек памяти
	var payload T

	err := json.NewDecoder(body).Decode(&payload)

	if err != nil {
		return nil, err
	}

	err = IsValid(payload)

	if err != nil {
		return nil, err
	}

	return &payload, nil
}
