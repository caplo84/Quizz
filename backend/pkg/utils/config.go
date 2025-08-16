package utils

import (
    "fmt"
    "os"
    "strconv"
    "strings"
)

// Config holds all application configuration
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Security SecurityConfig
    CORS     CORSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
    Port        string
    Environment string
    GinMode     string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
    URL      string
    Host     string
    Port     string
    Password string
    DB       int
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

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
    config := &Config{
        Server: ServerConfig{
            Port:        getEnvOrDefault("PORT", "8080"),
            Environment: getEnvOrDefault("ENVIRONMENT", "development"),
            GinMode:     getEnvOrDefault("GIN_MODE", "debug"),
        },
        Database: DatabaseConfig{
            Host:     getEnvOrDefault("DB_HOST", "localhost"),
            Port:     getEnvOrDefault("DB_PORT", "5432"),
            User:     getEnvOrDefault("DB_USER", "postgres"),
            Password: getEnvOrDefault("DB_PASSWORD", ""),
            DBName:   getEnvOrDefault("DB_NAME", "quiz_db"),
            SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            URL:      getEnvOrDefault("REDIS_URL", ""),
            Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
            Port:     getEnvOrDefault("REDIS_PORT", "6379"),
            Password: getEnvOrDefault("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
        Security: SecurityConfig{
            AdminAPIKey: os.Getenv("ADMIN_API_KEY"),
            JWTSecret:   os.Getenv("JWT_SECRET"),
        },
        CORS: CORSConfig{
            AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", ",", []string{"http://localhost:3000"}),
            AllowedMethods: getEnvAsSlice("ALLOWED_METHODS", ",", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
            AllowedHeaders: getEnvAsSlice("ALLOWED_HEADERS", ",", []string{"Origin", "Content-Type", "Authorization"}),
        },
    }
    
    if err := validateConfig(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
    if config.Security.AdminAPIKey == "" {
        return fmt.Errorf("ADMIN_API_KEY is required")
    }
    
    if config.Database.Password == "" && config.Server.Environment == "production" {
        return fmt.Errorf("DB_PASSWORD is required in production")
    }
    
    return nil
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