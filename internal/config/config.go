package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	AdminToken string
	HttpPort   string
	DBUrl      string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres_password"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "file_server"),
		getEnv("DB_SSLMODE", "disable"),
	)

	cfg := &Config{
		HttpPort:   getEnv("HTTP_PORT", "8080"),
		DBUrl:      dbUrl,
		JWTSecret:  getEnv("JWT_SECRET", "default_secret"),
		AdminToken: getEnv("ADMIN_TOKEN", "default_admin_token"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}
