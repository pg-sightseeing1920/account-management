package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadDB() DBConfig {
	cfg := DBConfig{
		Host:     get("DB_HOST", "localhost"),
		Port:     get("DB_PORT", "5432"),
		User:     get("DB_USER", "user"),
		Password: get("DB_PASSWORD", "password"),
		Name:     get("DB_NAME", "account_management"),
		SSLMode:  get("DB_SSLMODE", "disable"),
	}
	return cfg
}

func (c DBConfig) DSN() string {
	// postgres://user:pass@host:port/db?sslmode=disable
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

func get(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}