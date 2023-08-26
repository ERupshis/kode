package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Host       string
	LogLevel   string
	DbUser     string
	DbPassword string
	DbName     string
}

func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagAddress    = "a"
	flagDbSUser    = "db_user"
	flagDbSUserPwd = "db_password"
	flagDbName     = "db_name"
)

func checkFlags(config *Config) {
	flag.StringVar(&config.Host, flagAddress, "localhost:8080", "server endpoint")
	flag.StringVar(&config.DbUser, flagDbSUser, "postgres", "database super user name")
	flag.StringVar(&config.DbPassword, flagDbSUserPwd, "postgres", "database super user password")
	flag.StringVar(&config.DbName, flagDbName, "kodetest", "database name")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	Host       string `env:"ADDRESS"`
	DbUser     string `env:"DB_USER"`
	DbPassword string `env:"DB_PASSWORD"`
	DbName     string `env:"DB_NAME"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		panic(err)
	}

	setEnvToParamIfNeed(&config.Host, envs.Host)
	setEnvToParamIfNeed(&config.DbUser, envs.DbUser)
	setEnvToParamIfNeed(&config.DbPassword, envs.DbPassword)
	setEnvToParamIfNeed(&config.DbName, envs.DbName)
}
