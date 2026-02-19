package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
	Auth    AuthConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type StorageConfig struct {
	Type   StorageType
	Path   string
	Region string
	URL    string
}

type StorageType string

var (
	FSStorageType StorageType = "fs"
	S3StorageType StorageType = "s3"
)

type AuthConfig struct {
	Issuer string
	Secret string
	Exp    int
}

var (
	once     sync.Once
	instance *Config
	err      error
)

func GetConfig() (*Config, error) {
	once.Do(func() {
		instance, err = loadConfig()
	})

	return instance, err
}

func loadConfig() (*Config, error) {
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ошибка чтения конфига: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга: %w", err)
	}

	return &config, nil
}
