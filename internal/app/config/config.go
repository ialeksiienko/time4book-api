package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Facade struct {
	Env                     string `env:"ENV"`
	User                    string `env:"DB_USER"`
	Pass                    string `env:"DB_PASS"`
	Host                    string `env:"DB_HOST"`
	Port                    int    `env:"DB_PORT"`
	Name                    string `env:"DB_NAME"`
	JwtSecret               string `env:"JWT_SECRET"`
	JwtAccessTokenDuration  string `env:"JWT_ACCESS_TOKEN_DURATION"`
	JwtRefreshTokenDuration string `env:"JWT_REFRESH_TOKEN_DURATION"`
}

func MustLoad() *Facade {
	cfg, err := Load()
	if err != nil {
		panic("failed to read config: " + err.Error())
	}
	return cfg
}

func Load() (*Facade, error) {
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load .env: %w", err)
		}
	}

	var cfg Facade
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("read environment: %w", err)
	}

	{
		if strings.TrimSpace(cfg.Env) == "" {
			return nil, fmt.Errorf("ENV is required")
		}
		if strings.TrimSpace(cfg.User) == "" {
			return nil, fmt.Errorf("DB_USER is required")
		}
		if strings.TrimSpace(cfg.Pass) == "" {
			return nil, fmt.Errorf("DB_PASS is required")
		}
		if strings.TrimSpace(cfg.Host) == "" {
			return nil, fmt.Errorf("DB_HOST is required")
		}
		if strings.TrimSpace(cfg.Name) == "" {
			return nil, fmt.Errorf("DB_NAME is required")
		}
		if cfg.Port < 1 || cfg.Port > 65535 {
			return nil, fmt.Errorf("DB_PORT must be between 1 and 65535")
		}
		if strings.TrimSpace(cfg.JwtSecret) == "" {
			return nil, fmt.Errorf("JWT_SECRET is required")
		}
		if strings.TrimSpace(cfg.JwtAccessTokenDuration) == "" {
			return nil, fmt.Errorf("JWT_ACCESS_TOKEN_DURATION is required")
		}
		if strings.TrimSpace(cfg.JwtRefreshTokenDuration) == "" {
			return nil, fmt.Errorf("JWT_REFRESH_TOKEN_DURATION is required")
		}
		accessDuration, err := time.ParseDuration(cfg.JwtAccessTokenDuration)
		if err != nil || accessDuration <= 0 {
			return nil, fmt.Errorf("JWT_ACCESS_TOKEN_DURATION must be a positive duration")
		}
		refreshDuration, err := time.ParseDuration(cfg.JwtRefreshTokenDuration)
		if err != nil || refreshDuration <= 0 {
			return nil, fmt.Errorf("JWT_REFRESH_TOKEN_DURATION must be a positive duration")
		}
	}

	return &cfg, nil
}
