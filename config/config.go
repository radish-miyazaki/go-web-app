package config

import "github.com/caarlos0/env/v7"

type Config struct {
	Env        string `env:"APP_ENV" envDefault:"dev"`
	Port       int    `env:"PORT" envDefault:"80"`
	DBHost     string `env:"DB_HOST" envDefault:"127.0.0.1"`
	DBPort     int    `env:"DB_PORT" envDefault:"33060"`
	DBUser     string `env:"DB_USER" envDefault:"todo"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME" envDefault:"todo"`
	RedisHost  string `env:"REDIS_HOST" envDefault:"redis"`
	RedisPort  int    `env:"REDIS_PORT" envDefault:"36379"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
