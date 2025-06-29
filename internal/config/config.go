package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App       AppConfig      `yaml:"app"`
	Database  DatabaseConfig `yaml:"database"`
	Constants Constants      `yaml:"constants"`
}

type AppConfig struct {
	Server  ServerConfig  `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
	CORS    CORSConfig    `yaml:"cors"`
}

type ServerConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
	Timeout int    `yaml:"timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type DatabaseConfig struct {
	Driver         string               `yaml:"driver"`
	SQLite         SQLiteConfig         `yaml:"sqlite"`
	Postgres       PostgresConfig       `yaml:"postgres"`
	MySQL          MySQLConfig          `yaml:"mysql"`
	ConnectionPool ConnectionPoolConfig `yaml:"connection_pool"`
}

type SQLiteConfig struct {
	Path string `yaml:"path"`
}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	TimeZone string `yaml:"timezone"`
}

type MySQLConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"dbname"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse_time"`
	Loc       string `yaml:"loc"`
}

type ConnectionPoolConfig struct {
	MaxIdleConns    int `yaml:"max_idle_conns"`
	MaxOpenConns    int `yaml:"max_open_conns"`
	ConnMaxLifetime int `yaml:"conn_max_lifetime"`
}

type AuthConfig struct {
    JWTSecret       string `yaml:"jwt_secret"`
    TokenExpiration int    `yaml:"token_expiration"`
    PasswordCost    int    `yaml:"password_cost"`
}

type Constants struct {
    Pagination    PaginationConfig    `yaml:"pagination"`
    Validation    ValidationConfig    `yaml:"validation"`
    BusinessRules BusinessRulesConfig `yaml:"business_rules"`
    Auth          AuthConfig          `yaml:"auth"`
}

type PaginationConfig struct {
	DefaultPageSize int `yaml:"default_page_size"`
	MaxPageSize     int `yaml:"max_page_size"`
}

type ValidationConfig struct {
	MinPasswordLength int `yaml:"min_password_length"`
	MaxNameLength     int `yaml:"max_name_length"`
}

type BusinessRulesConfig struct {
	MaxProductsPerUser   int    `yaml:"max_products_per_user"`
	DefaultProductStatus string `yaml:"default_product_status"`
}

func Load() (*Config, error) {
	config := &Config{}

	// Load app config
	if err := loadYAMLFile("configs/app.yaml", &config.App); err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	// Load database config
	if err := loadYAMLFile("configs/database.yaml", &config.Database); err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load constants
	if err := loadYAMLFile("configs/constants.yaml", &config.Constants); err != nil {
		return nil, fmt.Errorf("failed to load constants: %w", err)
	}

	return config, nil
}

func loadYAMLFile(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.App.Server.Host, c.App.Server.Port)
}
