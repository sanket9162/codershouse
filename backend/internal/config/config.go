package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain    string
	Port      string
	SecretKey string
}

func LoadConfig() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.Port = os.Getenv("PORT")
	cfg.SecretKey = os.Getenv("SECRET_KEY")

	return &cfg

}
