package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Cache    CacheConfig    `mapstructure:"cache"`
	Features FeatureConfig  `mapstructure:"features"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port" env:"PORT" default:"8080"`
	Environment  string        `mapstructure:"environment" env:"ENVIRONMENT" default:"development"`
	GinMode      string        `mapstructure:"gin_mode" env:"GIN_MODE" default:"debug"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host" env:"DB_HOST" default:"localhost"`
	Port            string        `mapstructure:"port" env:"DB_PORT" default:"5432"`
	User            string        `mapstructure:"user" env:"DB_USER" default:"quiz_user"`
	Password        string        `mapstructure:"password" env:"DB_PASSWORD"`
	DBName          string        `mapstructure:"dbname" env:"DB_NAME" default:"quiz_db"`
	SSLMode         string        `mapstructure:"sslmode" env:"DB_SSLMODE" default:"disable"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host" env:"REDIS_HOST" default:"localhost"`
	Port         string `mapstructure:"port" env:"REDIS_PORT" default:"6379"`
	Password     string `mapstructure:"password" env:"REDIS_PASSWORD"`
	DB           int    `mapstructure:"db" env:"REDIS_DB" default:"0"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	URL          string `mapstructure:"url" env:"REDIS_URL"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level" env:"LOG_LEVEL" default:"info"`
	Format string `mapstructure:"format" env:"LOG_FORMAT" default:"json"`
	Output string `mapstructure:"output" env:"LOG_OUTPUT" default:"stdout"`
}

type SecurityConfig struct {
	AdminAPIKey  string `mapstructure:"admin_api_key"`
	JWTSecret    string `mapstructure:"jwt_secret"`
	RateLimitRPS int    `mapstructure:"rate_limit_rps"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type CacheConfig struct {
	TTLSeconds      int           `mapstructure:"ttl_seconds"`
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}

type FeatureConfig struct {
	EnableMetrics     bool `mapstructure:"enable_metrics"`
	EnableDebugRoutes bool `mapstructure:"enable_debug_routes"`
	EnableProfiling   bool `mapstructure:"enable_profiling"`
}

// LoadConfig loads configuration based on environment
func LoadConfig() (*Config, error) {
	// Get environment from ENV variable, default to development
	env := getEnvOrDefault("APP_ENV", "development")

	// Initialize viper
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath("/app/configs")

	// Enable environment variable substitution
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file when available; allow env/default-only mode in cloud deploys.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Expand environment variables in config values
	expandEnvVars()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	applyDefaults(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

// expandEnvVars replaces ${VAR} patterns with environment variables
func expandEnvVars() {
	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)
		if strings.Contains(value, "${") {
			expanded := os.ExpandEnv(value)
			viper.Set(key, expanded)
		}
	}
}

// validateConfig ensures required configuration is present
func validateConfig(config *Config) error {
	if config.Server.Environment == "production" {
		if config.Security.AdminAPIKey == "" || config.Security.AdminAPIKey == "dev_admin_key" {
			return fmt.Errorf("admin API key must be set for production")
		}
		if config.Security.JWTSecret == "" || config.Security.JWTSecret == "dev_jwt_secret" {
			return fmt.Errorf("JWT secret must be set for production")
		}
		if config.Database.SSLMode != "require" {
			return fmt.Errorf("SSL must be required for production database")
		}
	}

	return nil
}

func applyDefaults(config *Config) {
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Environment == "" {
		config.Server.Environment = "development"
	}
	if config.Server.GinMode == "" {
		config.Server.GinMode = "debug"
	}

	if config.Database.Port == "" {
		config.Database.Port = "5432"
	}
	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.User == "" {
		config.Database.User = "quiz_user"
	}
	if config.Database.DBName == "" {
		config.Database.DBName = "quiz_db"
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	if config.Redis.Host == "" {
		config.Redis.Host = "localhost"
	}
	if config.Redis.Port == "" {
		config.Redis.Port = "6379"
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "info"
	}
	if config.Logging.Format == "" {
		config.Logging.Format = "json"
	}
	if config.Logging.Output == "" {
		config.Logging.Output = "stdout"
	}
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
