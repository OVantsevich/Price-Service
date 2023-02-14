// Package config config
package config

import (
	"fmt"

	"github.com/caarlos0/env/v7"
)

// Config config
type Config struct {
	RedisHost     string `env:"REDIS_HOST,notEmpty" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT,notEmpty" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD,notEmpty" envDefault:"redis"`
	StreamName    string `env:"STREAM_NAME,notEmpty" envDefault:"prices"`
	GroupName     string `env:"GROUP_NAME,notEmpty" envDefault:"priceGroup"`
	Port          string `env:"PORT,notEmpty" envDefault:"10000"`
}

// NewConfig parsing config from environment
func NewConfig() (*Config, error) {
	config := &Config{}

	err := env.Parse(config)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - Parse:%w", err)
	}

	return config, nil
}
