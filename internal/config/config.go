package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Host     string
	LogLevel string
}

func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagAddress = "a"
)

func checkFlags(config *Config) {
	flag.StringVar(&config.Host, flagAddress, "localhost:8080", "server endpoint")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	Host string `env:"ADDRESS"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		panic(err)
	}

	setEnvToParamIfNeed(&config.Host, envs.Host)
}
