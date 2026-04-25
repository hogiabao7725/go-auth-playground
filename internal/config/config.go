package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	Logger LoggerConfig
	DB     DBConfig
	JWT    JWTConfig
}

type LoggerConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	Format string `env:"LOG_FORMAT" envDefault:"text"`
	Pretty bool   `env:"LOG_PRETTY" envDefault:"true"`
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
	Env  string `env:"ENV" envDefault:"development"`
	Name string `env:"APP_NAME" envDefault:"go-auth-playground"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
	MaxConns int    `env:"DB_MAX_CONNS" envDefault:"25"`
	MinConns int    `env:"DB_MIN_CONNS" envDefault:"5"`

	ConnectTimeout time.Duration `env:"DB_CONNECT_TIMEOUT" envDefault:"10s"`
	ConnLifetime   time.Duration `env:"DB_CONN_LIFETIME" envDefault:"1h"`
	ConnIdleTime   time.Duration `env:"DB_CONN_IDLE_TIME" envDefault:"30m"`
}

type JWTConfig struct {
	AccessSecret  string        `env:"JWT_ACCESS_SECRET,required"`
	RefreshSecret string        `env:"JWT_REFRESH_SECRET,required"`
	AccessTTL     time.Duration `env:"JWT_ACCESS_TTL" envDefault:"15m"`
	RefreshTTL    time.Duration `env:"JWT_REFRESH_TTL" envDefault:"168h"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config load failed:\n %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	var errs []string
	const minJWTSecretLength = 32

	if len(c.JWT.AccessSecret) < minJWTSecretLength {
		errs = append(errs, "jwt: JWT_ACCESS_SECRET must be at least 32 characters")
	}
	if len(c.JWT.RefreshSecret) < minJWTSecretLength {
		errs = append(errs, "jwt: JWT_REFRESH_SECRET must be at least 32 characters")
	}
	if c.DB.MaxConns < c.DB.MinConns {
		errs = append(errs, "database: DB_MAX_CONNS must be >= DB_MIN_CONNS")
	}
	if c.DB.ConnectTimeout <= 0 {
		errs = append(errs, "database: DB_CONNECT_TIMEOUT must be > 0")
	}
	if c.DB.ConnLifetime <= 0 {
		errs = append(errs, "database: DB_CONN_LIFETIME must be > 0")
	}
	if c.DB.ConnIdleTime <= 0 {
		errs = append(errs, "database: DB_CONN_IDLE_TIME must be > 0")
	}

	if c.Logger.Format != "text" && c.Logger.Format != "json" {
		errs = append(errs, "logger: LOG_FORMAT must be 'text' or 'json'")
	}

	if len(errs) > 0 {
		return fmt.Errorf("config validation failed:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func (d *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}
