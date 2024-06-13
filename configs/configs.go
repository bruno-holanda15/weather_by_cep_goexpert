package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type conf struct {
	WeatherApiToken string `mapstructure:"WEATHER_API_TOKEN"`
}

// fullcycle way of loading envs
// but, I don't like the idea of when I want to
// load a specific env, I need to LoadConfig then
// put conf.ENV_I_WANT for example
func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}

// I think it's easier and direct that way
func Env(envAndDefault ...string) string {
	if err := godotenv.Load(); err != nil {
		panic("NOT ABLE TO LOAD ENVS")
	}

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
