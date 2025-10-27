package services

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	PostgresName     string
}

func LoadConfig() *Config {
	// Load .env file if it exists (ignore error if file doesn't exist in production)
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warn("No .env file found, using environment variables")
	}

	return &Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresName:     os.Getenv("POSTGRES_DB"),
	}
}

func (c *Config) PostgresURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.PostgresUser, c.PostgresPassword, c.PostgresHost, c.PostgresPort, c.PostgresName,
	)
}
