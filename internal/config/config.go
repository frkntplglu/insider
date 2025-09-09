package config

import (
	"os"
	"strconv"
	"time"

	"github.com/frkntplglu/insider/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Ticker   TickerConfig
	SMS      SMSConfig
}

type AppConfig struct {
	Name    string
	Version string
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Database int
	Password string
	Key      string
}

type TickerConfig struct {
	Period time.Duration
}

type SMSConfig struct {
	Host string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("No .env file found, using default values")
	}
	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "Insider Auto Message Sender"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "insider"),
			Password: getEnv("DB_PASSWORD", "insider"),
			Database: getEnv("DB_NAME", "insiderdb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Database: getEnvAsInt("REDIS_DB", 0),
			Password: getEnv("REDIS_PASSWORD", ""),
			Key:      getEnv("REDIS_KEY", "sent:messages"),
		},
		Ticker: TickerConfig{
			Period: getEnvAsDuration("TICKER_PERIOD", 10*time.Minute),
		},
		SMS: SMSConfig{
			Host: getEnv("SMS_HOST", "https://webhook.site/92566df1-0f3b-41f9-89a0-dae31e91c565"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
