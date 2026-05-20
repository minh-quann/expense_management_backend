package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	AppPort  string
	GinMode  string
	Database DatabaseConfig
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

// DSN returns the PostgreSQL connection string
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode, d.Timezone,
	)
}

// LoadConfig reads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "expense_user"),
			Password: getEnv("DB_PASSWORD", "expense_secret"),
			DBName:   getEnv("DB_NAME", "expense_management"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "Asia/Ho_Chi_Minh"),
		},
	}
}

// getEnv reads an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
