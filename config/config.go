package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db Db

	AppPort string
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbUser := os.Getenv("DATABASE_USER")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")

	appPort := os.Getenv("PORT")

	dbConfig := Db{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DbName:   dbName,
	}

	return &Config{
		Db:      dbConfig,
		AppPort: appPort,
	}
}
