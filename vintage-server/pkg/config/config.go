package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config adalah struct yang akan menampung semua konfigurasi aplikasi.
// Tag `mapstructure` digunakan oleh Viper untuk memetakan nama variabel.
type Config struct {
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          int    `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	DBName          string `mapstructure:"DB_NAME"`
	DBSSLMode       string `mapstructure:"DB_SSLMODE"`
	UserServicePort int    `mapstructure:"USER_SERVICE_PORT"`
	JWTSecretKey    string `mapstructure:"JWT_SECRET_KEY"`
	CloudinaryURL   string `mapstructure:"CLOUDINARY_URL"`
}

// DSN (Data Source Name) mengembalikan connection string untuk database.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// LoadConfig memuat konfigurasi dari file .env dan environment variables.
func LoadConfig() (config Config, err error) {
	// Memuat file .env jika ada (berguna untuk development lokal)
	// Di lingkungan production, kita akan set environment variables langsung.
	err = godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file")
	}

	// Memberi tahu Viper untuk membaca environment variables
	viper.AutomaticEnv()

	// Mengikat environment variables ke field di struct Config
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("USER_SERVICE_PORT")
	viper.BindEnv("JWT_SECRET_KEY")
	viper.BindEnv("CLOUDINARY_URL")

	// Unmarshal semua konfigurasi yang ditemukan ke dalam struct Config
	err = viper.Unmarshal(&config)
	return
}
