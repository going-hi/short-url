package utils

import (
	"math/rand/v2"
)

var Chars = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateCode() string {
	code := make([]rune, 6)

	for index, _ := range code {
		char := rand.IntN(len(Chars) - 1)
		code[index] = Chars[char]
	}
	return string(code)
}
