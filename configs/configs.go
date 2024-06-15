package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Loader struct{}

// I think it's easier and direct that way
func Env(envAndDefault ...string) string {
	if len(envAndDefault) == 0 {
		return ""
	}
	key := envAndDefault[0]
	env := os.Getenv(key)

	if env != "" {
		return env
	}

	defaultValue := ""
	if len(envAndDefault) > 1 && envAndDefault[1] != "" {
		defaultValue = envAndDefault[1]
	}

	if defaultValue == "" {
		panic("unable to read environment variable: " + key)
	}

	return defaultValue
}

func (c *Loader) LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("unable to load environment vars")
	}
}

