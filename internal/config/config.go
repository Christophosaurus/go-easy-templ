package config

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var configFile = "config.yml"

type Config struct {
	Server Server `yaml:"server"`
	DB     DB     `yaml:"database"`
}

type Server struct {
	Hostname     string `yaml:"host" envconfig:"SERVER_HOSTNAME"`
	Port         string `yaml:"port" envconfig:"SERVER_PORT"`
	Env          string `yaml:"env" envconfig:"SERVER_ENV"`
	LoggingLevel string `yaml:"logging_level" envconfig:"SERVER_LOGGING_LEVEL"`
}

type DB struct {
	DSN          string `yaml:"dsn" envconfig:"DB_DSN"`
	MaxOpenConns int    `yaml:"max_open_conns" envconfig:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns int    `yaml:"max_idle_conns" envconfig:"DB_MAX_IDLE_CONNS"`
	MaxIdleTime  string `yaml:"max_idle_time" envconfig:"DB_MAX_IDLE_TIME"`
}

func InitConfig() (*Config, error) {
	config := &Config{}

	// Read from config.yml
	err := ReadFile(config, configFile)
	if err != nil {

		return nil, err
	}
	// Read from env variables and overwrite
	err = ReadEnv(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ReadFile(config *Config, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

func ReadEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}
	return nil
}
