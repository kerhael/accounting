package config

import (
	"fmt"
	"os"
	"strings"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	Database DatabaseConfig
}

func Load() (*Config, error) {
	var cfgErr []string
	if os.Getenv("DB_HOST") == "" {
		cfgErr = append(cfgErr, "DB_HOST")
	}
	if os.Getenv("DB_PORT") == "" {
		cfgErr = append(cfgErr, "DB_PORT")
	}
	if os.Getenv("DB_USER") == "" {
		cfgErr = append(cfgErr, "DB_USER")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		cfgErr = append(cfgErr, "DB_PASSWORD")
	}
	if os.Getenv("DB_NAME") == "" {
		cfgErr = append(cfgErr, "DB_NAME")
	}
	if os.Getenv("DB_SSLMODE") == "" {
		cfgErr = append(cfgErr, "DB_SSLMODE")
	}
	if len(cfgErr) > 0 {
		return nil, fmt.Errorf("missing %s", strings.Join(cfgErr, ","))
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	return cfg, nil
}
