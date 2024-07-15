package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// struct tag позволяют другому коду на go оценивать структуры и извлекать значения

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-requred:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// must честное именование, паникуем, если ошибка
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

// fetchConfigPath fetches config path from command line flag or env var
// Priority: flag > env > default
// Default value - empty string
// CONFIG_PATH=./path/to/config.yaml myApp
// myApp --config=./path/to/config.yaml
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	// подсказка для командной строки
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
