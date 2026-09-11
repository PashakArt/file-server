package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Config struct {
	JWTSecret             string
	AdminToken            string
	HttpPort              string
	DBUrl                 string
	Redis                 RedisConfig
	TokenTTL              time.Duration
	MaxUploadBodySizeByte int
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
		HttpPort:              getEnv("HTTP_PORT", "8080"),
		DBUrl:                 dbUrl,
		JWTSecret:             getEnv("JWT_SECRET", "default_secret"),
		AdminToken:            getEnv("ADMIN_TOKEN", "default_admin_token"),
		TokenTTL:              getEnvAsDuration("TOKEN_TTL", 24*time.Hour),
		MaxUploadBodySizeByte: getEnvAsInt("MAX_UPLOAD_FILE_SIZE_MB", 32),
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		}}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr, ok := os.LookupEnv(key)
	if ok {
		v, _ := strconv.Atoi(valueStr)
		return v
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr, ok := os.LookupEnv(key)
	if !ok || valueStr == "" {
		return defaultValue
	}

	val, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return val
}
