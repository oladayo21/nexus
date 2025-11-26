package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env         string `env:"ENV" envDefault:"development"`
	Port        int    `env:"PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	RedisURL    string `env:"REDIS_URL,required"`
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
