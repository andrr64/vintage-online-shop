package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTSecretAdmin string
	DevKey         string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTSecretAdmin: os.Getenv("JWT_SECRET_ADMIN"),
		DevKey:         os.Getenv("DEV_KEY"),
	}

	if AppConfig.DBHost == "" || AppConfig.DBUser == "" {
		log.Fatal("Missing required environment variables")
	}
}
