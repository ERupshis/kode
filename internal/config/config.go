package config

import (
	"flag"
	"log"

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
	flagAddress  = "a"
	flagLogLevel = "l"
)

func checkFlags(config *Config) {
	flag.StringVar(&config.Host, flagAddress, "localhost:8080", "server endpoint")
	flag.StringVar(&config.LogLevel, flagLogLevel, "Info", "log level")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	Host     string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	setEnvToParamIfNeed(&config.Host, envs.Host)
	setEnvToParamIfNeed(&config.LogLevel, envs.LogLevel)
}
