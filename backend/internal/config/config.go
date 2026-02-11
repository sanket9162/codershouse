package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sanket9162/codershouse/internal/utils"
)

type Config struct {
	Domain      string
	Port        string
	SecretKey   string
	TwilioSID   string
	TwilioToken string
	TwilioPhone string
	DBUri       string
	DBName      string
	Auth        utils.Auth
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
	cfg.TwilioPhone = os.Getenv("TWILIO_PHONE_NUMBER")
	cfg.DBUri = os.Getenv("DB_URL")
	cfg.DBName = os.Getenv("DB_NAME")

	//Auth Configuration
	cfg.Auth = utils.Auth{
		Issuer:        os.Getenv("JWT_ISSUER"),
		Audience:      os.Getenv("JWT_AUDIENCE"),
		Secret:        os.Getenv("JWT_SECRET"),
		TokenExpiry:   time.Hour * 1,
		RefreshExpiry: time.Hour * 24 * 7,
		CookieDomain:  os.Getenv("COOKIE_DOMAIN"),
		CookiePath:    os.Getenv("COOKIE_PATH"),
		CookieName:    os.Getenv("COOKIE_NAME"),
	}

	return &cfg

}
