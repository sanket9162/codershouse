package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain      string
	Port        string
	SecretKey   string
	TwilioSID   string
	TwilioToken string
	TwilioPhone string
}

func LoadConfig() *Config {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg.Domain = os.Getenv("DOMAIN")
	cfg.Port = os.Getenv("PORT")
	cfg.SecretKey = os.Getenv("KEY")
	cfg.TwilioSID = os.Getenv("TWILIO_SID")
	cfg.TwilioToken = os.Getenv("TWILIO_TOKEN")
	cfg.TwilioPhone = os.Getenv("TWILIO_PHONE")

	return &cfg

}
