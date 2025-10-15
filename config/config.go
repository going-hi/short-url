package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db Db

	SecretKey string
	AppPort   string
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

// Можно добавить проверку обязательный параметров
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("PORT")
	secretKey := os.Getenv("SECRET_KEY_JWT")

	dbConfig := Db{ // вот так логичнее выглядит и меньше, и не создаются лишенне переменные значит выигрываем по памяти
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DbName:   os.Getenv("DATABASE_NAME"),
	}
	return &Config{
		Db:        dbConfig,
		AppPort:   appPort,
		SecretKey: secretKey,
	}
}
