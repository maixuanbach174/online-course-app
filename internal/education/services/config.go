package services

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresName     string
}

func LoadConfig() *Config {
	return &Config{
		PostgresUser:     os.Getenv("DB_USER"),
		PostgresPassword: os.Getenv("DB_PASSWORD"),
		PostgresHost:     os.Getenv("DB_HOST"),
		PostgresPort:     os.Getenv("DB_PORT"),
		PostgresName:     os.Getenv("DB_NAME"),
	}
}

func (c *Config) PostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresName,
	)
}
