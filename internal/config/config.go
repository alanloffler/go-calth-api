package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppDomain        string
	DatabaseURL      string
	Port             string
	JwtSecret        string
	JwtRefreshSecret string
	JwtAccessExpiry  string
	JwtRefreshExpiry string
	CookieDomain     string
	CookieSecure     bool
	CorsOrigin       string
	RedisAddr        string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Failed to load .env file")
	}

	var config *Config = &Config{
		AppDomain:        os.Getenv("APP_DOMAIN"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		Port:             os.Getenv("PORT"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		JwtRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		JwtAccessExpiry:  os.Getenv("JWT_ACCESS_EXPIRY"),
		JwtRefreshExpiry: os.Getenv("JWT_REFRESH_EXPIRY"),
		CookieDomain:     os.Getenv("COOKIE_DOMAIN"),
		CookieSecure:     os.Getenv("COOKIE_SECURE") == "true",
		CorsOrigin:       os.Getenv("CORS_ORIGIN"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
	}

	return config, nil
}
