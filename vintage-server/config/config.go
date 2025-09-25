package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	ServerPort     string
}

var AppConfig Config

func LoadConfig() {

	err := godotenv.Load()
	log.Println("Loading .env file")
	if err != nil {
		log.Println(".env file not found, reading environment variables directly")
	}

	AppConfig = Config{
		DBHost:         os.Getenv("DB_HOST"),
		DBPort:         os.Getenv("DB_PORT"),
		DBUser:         os.Getenv("DB_USER"),
		DBPassword:     os.Getenv("DB_PASSWORD"),
		DBName:         os.Getenv("DB_NAME"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTSecretAdmin: os.Getenv("JWT_SECRET_ADMIN"),
		DevKey:         os.Getenv("DEV_KEY"),
		ServerPort:     os.Getenv("SERVER_PORT"),
	}

	if AppConfig.DBHost == "" || AppConfig.DBUser == "" {
		log.Fatal("Missing required environment variables")
	}
}
