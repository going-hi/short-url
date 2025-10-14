package utils

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

func IsValid(payload any) error {
	var validate = validator.New()
	err := validate.Struct(payload)
	return err
}

func SendJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func GetBody[T any](body io.ReadCloser) (*T, error) {
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
