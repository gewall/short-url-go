package pkg

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env   string
	DBURL string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfig() error {
	env := os.Getenv("FOO_ENV")
	if env == "" {
		env = "development"
	}

	if env == "development" {
		_ = godotenv.Load(".env.development")
	} else {
		_ = godotenv.Load()
	}

	c.Env = os.Getenv("FOO_ENV")
	c.DBURL = os.Getenv("DATABASE_URL")

	return nil
}
