package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Facade struct {
	Env string `env:"ENV"`
	//
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Host string `env:"DB_HOST"`
	Port int    `env:"DB_PORT"`
	Name string `env:"DB_NAME"`
	//
	JwtSecret               string `env:"JWT_SECRET"`
	JwtAccessTokenDuration  string `env:"JWT_ACCESS_TOKEN_DURATION"`
	JwtRefreshTokenDuration string `env:"JWT_REFRESH_TOKEN_DURATION"`
}

func MustLoad() *Facade {
	var cfg Facade

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
