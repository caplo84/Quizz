package utils

import (
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Logging  LoggingConfig  `mapstructure:"logging"`
    Security SecurityConfig `mapstructure:"security"`
    CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
    Port        string `mapstructure:"port" env:"PORT" default:"8080"`
    Environment string `mapstructure:"environment" env:"ENVIRONMENT" default:"development"`
    GinMode  string          `mapstructure:"gin_mode" env:"GIN_MODE" default:"debug"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host" env:"DB_HOST" default:"localhost"`
    Port     string `mapstructure:"port" env:"DB_PORT" default:"5432"`
    User     string `mapstructure:"user" env:"DB_USER" default:"quiz_user"`
    Password string `mapstructure:"password" env:"DB_PASSWORD"`
    DBName   string `mapstructure:"dbname" env:"DB_NAME" default:"quiz_db"`
    SSLMode  string `mapstructure:"sslmode" env:"DB_SSLMODE" default:"disable"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host" env:"REDIS_HOST" default:"localhost"`
    Port     string `mapstructure:"port" env:"REDIS_PORT" default:"6379"`
    Password string `mapstructure:"password" env:"REDIS_PASSWORD"`
    DB       int    `mapstructure:"db" env:"REDIS_DB" default:"0"`
    URL string `mapstructure:"url" env:"REDIS_URL"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
    AdminAPIKey string
    JWTSecret   string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
    AllowedOrigins []string
    AllowedMethods []string
    AllowedHeaders []string
}

type LoggingConfig struct {
    Level  string `mapstructure:"level" env:"LOG_LEVEL" default:"info"`
    Format string `mapstructure:"format" env:"LOG_FORMAT" default:"json"`
    Output string `mapstructure:"output" env:"LOG_OUTPUT" default:"stdout"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    redisDB, _ := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))
    
    config := &Config{
        Server: ServerConfig{
            Port:        getEnvOrDefault("PORT", "8080"),
            Environment: getEnvOrDefault("ENVIRONMENT", "development"),
            GinMode:     getEnvOrDefault("GIN_MODE", "debug"),
        },
        Database: DatabaseConfig{
            Host:     getEnvOrDefault("DB_HOST", "localhost"),
            Port:     getEnvOrDefault("DB_PORT", "5432"),
            User:     getEnvOrDefault("DB_USER", "quiz_user"),
            Password: getEnvOrDefault("DB_PASSWORD", ""),
            DBName:   getEnvOrDefault("DB_NAME", "quiz_db"),
            SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
            Port:     getEnvOrDefault("REDIS_PORT", "6379"),
            Password: getEnvOrDefault("REDIS_PASSWORD", ""),
            DB:       redisDB,
            URL:      getEnvOrDefault("REDIS_URL", ""),
        },
        Security: SecurityConfig{
            AdminAPIKey: getEnvOrDefault("ADMIN_API_KEY", "admin-key"),
            JWTSecret:   getEnvOrDefault("JWT_SECRET", "jwt-secret"),
        },
        CORS: CORSConfig{
            AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", ",", []string{"http://localhost:3000"}),
            AllowedMethods: getEnvAsSlice("ALLOWED_METHODS", ",", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
            AllowedHeaders: getEnvAsSlice("ALLOWED_HEADERS", ",", []string{"Origin", "Content-Type", "Authorization"}),
        },
        Logging: LoggingConfig{
            Level:  getEnvOrDefault("LOG_LEVEL", "info"),
            Format: getEnvOrDefault("LOG_FORMAT", "json"),
            Output: getEnvOrDefault("LOG_OUTPUT", "stdout"),
        },
    }
    
    return config, nil
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
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

func getEnvAsSlice(key, separator string, defaultValue []string) []string {
    if value := os.Getenv(key); value != "" {
        return strings.Split(value, separator)
    }
    return defaultValue
}