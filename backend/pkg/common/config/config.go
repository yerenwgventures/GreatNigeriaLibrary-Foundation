package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

// Config represents application configuration
type Config struct {
	Server      ServerConfig      `json:"server" yaml:"server"`
	Database    DatabaseConfig    `json:"database" yaml:"database"`
	Redis       RedisConfig       `json:"redis" yaml:"redis"`
	Auth        AuthConfig        `json:"auth" yaml:"auth"`
	OAuth       OAuthConfig       `json:"oauth" yaml:"oauth"`
	Email       EmailConfig       `json:"email" yaml:"email"`
	Storage     StorageConfig     `json:"storage" yaml:"storage"`
	Logging     LoggingConfig     `json:"logging" yaml:"logging"`
	Services    ServicesConfig    `json:"services" yaml:"services"`
	Features    FeaturesConfig    `json:"features" yaml:"features"`
	RateLimit   RateLimitConfig   `json:"rate_limiting" yaml:"rate_limiting"`
	CORS        CORSConfig        `json:"cors" yaml:"cors"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	Environment  string        `json:"environment" yaml:"environment"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string `json:"host"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	SSLMode         string `json:"ssl_mode"`
	MaxOpenConns    int    `json:"max_open_conns"`
	MaxIdleConns    int    `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// RedisConfig represents enhanced Redis configuration
type RedisConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Password     string        `json:"password" yaml:"password"`
	Database     int           `json:"database" yaml:"database"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns" yaml:"min_idle_conns"`
	MaxRetries   int           `json:"max_retries" yaml:"max_retries"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	Enabled      bool          `json:"enabled" yaml:"enabled"`
}

// AuthConfig represents enhanced authentication configuration
type AuthConfig struct {
	JWTSecret                   string        `json:"jwt_secret" yaml:"jwt_secret"`
	JWTIssuer                   string        `json:"jwt_issuer" yaml:"jwt_issuer"`
	AccessTokenExpiration       time.Duration `json:"access_token_expiration" yaml:"access_token_expiration"`
	RefreshTokenExpiration      time.Duration `json:"refresh_token_expiration" yaml:"refresh_token_expiration"`
	PasswordResetExpiration     time.Duration `json:"password_reset_expiration" yaml:"password_reset_expiration"`
	EmailVerificationExpiration time.Duration `json:"email_verification_expiration" yaml:"email_verification_expiration"`
	SessionExpiration           time.Duration `json:"session_expiration" yaml:"session_expiration"`
	MaxLoginAttempts            int           `json:"max_login_attempts" yaml:"max_login_attempts"`
	LockoutDuration             time.Duration `json:"lockout_duration" yaml:"lockout_duration"`
	EnableTokenRevocation       bool          `json:"enable_token_revocation" yaml:"enable_token_revocation"`
	EnableSessionTracking       bool          `json:"enable_session_tracking" yaml:"enable_session_tracking"`
	RequireEmailVerification    bool          `json:"require_email_verification" yaml:"require_email_verification"`
	Enable2FA                   bool          `json:"enable_2fa" yaml:"enable_2fa"`
	TokenSecurityChecks         bool          `json:"token_security_checks" yaml:"token_security_checks"`
}

// OAuthConfig represents OAuth configuration
type OAuthConfig struct {
	Google GoogleOAuthConfig `json:"google"`
}

// GoogleOAuthConfig represents Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}

// EmailConfig represents email configuration
type EmailConfig struct {
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	FromEmail    string `json:"from_email"`
	FromName     string `json:"from_name"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	Type      string `json:"type"` // local, s3, gcs
	LocalPath string `json:"local_path"`
	S3Config  S3Config `json:"s3"`
}

// S3Config represents S3 storage configuration
type S3Config struct {
	Region    string `json:"region"`
	Bucket    string `json:"bucket"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"` // json, text
	Output string `json:"output" yaml:"output"` // stdout, file
	File   string `json:"file" yaml:"file"`
}

// ServicesConfig represents microservices configuration
type ServicesConfig struct {
	AuthService       ServiceConfig `json:"auth_service" yaml:"auth_service"`
	ContentService    ServiceConfig `json:"content_service" yaml:"content_service"`
	DiscussionService ServiceConfig `json:"discussion_service" yaml:"discussion_service"`
	APIGateway        ServiceConfig `json:"api_gateway" yaml:"api_gateway"`
}

// ServiceConfig represents individual service configuration
type ServiceConfig struct {
	Port int `json:"port" yaml:"port"`
}

// FeaturesConfig represents feature flags
type FeaturesConfig struct {
	EnableRegistration     bool `json:"enable_registration" yaml:"enable_registration"`
	EnableOAuth           bool `json:"enable_oauth" yaml:"enable_oauth"`
	EnableEmailVerification bool `json:"enable_email_verification" yaml:"enable_email_verification"`
	EnableTwoFactorAuth   bool `json:"enable_two_factor_auth" yaml:"enable_two_factor_auth"`
	EnableFileUploads     bool `json:"enable_file_uploads" yaml:"enable_file_uploads"`
	EnableDiscussions     bool `json:"enable_discussions" yaml:"enable_discussions"`
	MaintenanceMode       bool `json:"maintenance_mode" yaml:"maintenance_mode"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool `json:"enabled" yaml:"enabled"`
	RequestsPerMinute int  `json:"requests_per_minute" yaml:"requests_per_minute"`
	BurstSize         int  `json:"burst_size" yaml:"burst_size"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvAsInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 120*time.Second),
			Environment:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvAsInt("DB_PORT", 5432),
			Username:        getEnv("DB_USERNAME", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_DATABASE", "great_nigeria_library"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			Database:     getEnvAsInt("REDIS_DATABASE", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
			MaxRetries:   getEnvAsInt("REDIS_MAX_RETRIES", 3),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
			Enabled:      getEnvAsBool("REDIS_ENABLED", true),
		},
		Auth: AuthConfig{
			JWTSecret:                   getEnv("JWT_SECRET", "your-secret-key"),
			JWTIssuer:                   getEnv("JWT_ISSUER", "great-nigeria-library"),
			AccessTokenExpiration:       getEnvAsDuration("ACCESS_TOKEN_EXPIRATION", 15*time.Minute),
			RefreshTokenExpiration:      getEnvAsDuration("REFRESH_TOKEN_EXPIRATION", 7*24*time.Hour),
			PasswordResetExpiration:     getEnvAsDuration("PASSWORD_RESET_EXPIRATION", 1*time.Hour),
			EmailVerificationExpiration: getEnvAsDuration("EMAIL_VERIFICATION_EXPIRATION", 24*time.Hour),
			SessionExpiration:           getEnvAsDuration("SESSION_EXPIRATION", 30*24*time.Hour),
			MaxLoginAttempts:            getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
			LockoutDuration:             getEnvAsDuration("LOCKOUT_DURATION", 15*time.Minute),
			EnableTokenRevocation:       getEnvAsBool("ENABLE_TOKEN_REVOCATION", true),
			EnableSessionTracking:       getEnvAsBool("ENABLE_SESSION_TRACKING", true),
			RequireEmailVerification:    getEnvAsBool("REQUIRE_EMAIL_VERIFICATION", false),
			Enable2FA:                   getEnvAsBool("ENABLE_2FA", false),
			TokenSecurityChecks:         getEnvAsBool("TOKEN_SECURITY_CHECKS", true),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			},
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "localhost"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@greatnigeria.com"),
			FromName:     getEnv("FROM_NAME", "Great Nigeria Library"),
		},
		Storage: StorageConfig{
			Type:      getEnv("STORAGE_TYPE", "local"),
			LocalPath: getEnv("STORAGE_LOCAL_PATH", "./uploads"),
			S3Config: S3Config{
				Region:    getEnv("S3_REGION", "us-east-1"),
				Bucket:    getEnv("S3_BUCKET", ""),
				AccessKey: getEnv("S3_ACCESS_KEY", ""),
				SecretKey: getEnv("S3_SECRET_KEY", ""),
			},
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			Output: getEnv("LOG_OUTPUT", "stdout"),
			File:   getEnv("LOG_FILE", "app.log"),
		},
	}

	return config, nil
}

// GetOAuthConfig returns OAuth configuration for a provider
func (c *Config) GetOAuthConfig(provider string) *oauth2.Config {
	switch provider {
	case "google":
		return &oauth2.Config{
			ClientID:     c.OAuth.Google.ClientID,
			ClientSecret: c.OAuth.Google.ClientSecret,
			RedirectURL:  c.OAuth.Google.RedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		}
	default:
		return nil
	}
}

// GetDatabaseDSN returns database connection string
func (c *Config) GetDatabaseDSN() string {
	return "host=" + c.Database.Host +
		" port=" + strconv.Itoa(c.Database.Port) +
		" user=" + c.Database.Username +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.Database +
		" sslmode=" + c.Database.SSLMode
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port)
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// Helper functions
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// LoadFromYAML loads configuration from YAML file with environment variable overrides
func LoadFromYAML(filename string) (*Config, error) {
	// First try to load from YAML file
	config := &Config{}

	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	}

	// Override with environment variables (environment takes precedence)
	envConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment config: %w", err)
	}

	// Merge configurations (env overrides YAML)
	mergeConfigs(config, envConfig)

	return config, nil
}

// mergeConfigs merges environment config into YAML config
func mergeConfigs(yamlConfig, envConfig *Config) {
	// Server config
	if envConfig.Server.Host != "0.0.0.0" || yamlConfig.Server.Host == "" {
		yamlConfig.Server.Host = envConfig.Server.Host
	}
	if envConfig.Server.Port != 8080 || yamlConfig.Server.Port == 0 {
		yamlConfig.Server.Port = envConfig.Server.Port
	}
	if envConfig.Server.Environment != "development" || yamlConfig.Server.Environment == "" {
		yamlConfig.Server.Environment = envConfig.Server.Environment
	}

	// Database config - environment variables always override YAML
	if os.Getenv("DB_HOST") != "" {
		yamlConfig.Database.Host = envConfig.Database.Host
	}
	if os.Getenv("DB_PORT") != "" {
		yamlConfig.Database.Port = envConfig.Database.Port
	}
	if os.Getenv("DB_USERNAME") != "" {
		yamlConfig.Database.Username = envConfig.Database.Username
	}
	if os.Getenv("DB_PASSWORD") != "" {
		yamlConfig.Database.Password = envConfig.Database.Password
	}
	if os.Getenv("DB_DATABASE") != "" {
		yamlConfig.Database.Database = envConfig.Database.Database
	}

	// Auth config - environment variables always override YAML
	if os.Getenv("JWT_SECRET") != "" {
		yamlConfig.Auth.JWTSecret = envConfig.Auth.JWTSecret
	}

	// Redis config - environment variables always override YAML
	if os.Getenv("REDIS_HOST") != "" {
		yamlConfig.Redis.Host = envConfig.Redis.Host
	}
	if os.Getenv("REDIS_PORT") != "" {
		yamlConfig.Redis.Port = envConfig.Redis.Port
	}
	if os.Getenv("REDIS_PASSWORD") != "" {
		yamlConfig.Redis.Password = envConfig.Redis.Password
	}
}
