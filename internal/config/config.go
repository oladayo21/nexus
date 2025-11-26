package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env         string `env:"ENV" envDefault:"development"`
	Port        int    `env:"PORT" envDefault:"8080"`
	AppSecret   string `env:"APP_SECRET,required"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	RedisAddr   string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}
