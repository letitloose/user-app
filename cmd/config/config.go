package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Db DBConfig
}

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
}

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

func (config *Config) ReadConfig(fileName string) error {
	log.Printf("reading config: %s\n", fileName)
	configFileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return errors.New(fmt.Sprintf("error opening config file: %s\n", err))
	}

	err = yaml.Unmarshal(configFileBytes, config)
	if err != nil {
		return errors.New(fmt.Sprintf("error unmarshalling config file:%s\n", err))
	}

	log.Printf("Config Values: %v\n", config)

	return nil
}
