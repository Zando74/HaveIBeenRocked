package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		once     sync.Once
		instance *Config
		App      `yaml:"app"`
	}

	App struct {
		Name     string `yaml:"name" `
		Version  string `yaml:"version"`
		Database `yaml:"database"`
		Log      `yaml:"log"`
		Grpc     `yaml:"grpc"`
		Http     `yaml:"http"`
	}

	Database struct {
		Host         string `yaml:"host" env:"DB_HOST"`
		Port         string `yaml:"port" env:"DB_PORT"`
		User         string `yaml:"user" env:"DB_USER"`
		Password     string `yaml:"password" env:"DB_PASSWORD"`
		RetryTimeout int    `yaml:"retry_timeout" env:"DB_RETRY_TIMEOUT"`
		MaxIdleConns int    `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
		MaxOpenConns int    `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	Grpc struct {
		Port int `yaml:"port" env:"GRPC_PORT"`
	}

	Http struct {
		Port int `yaml:"port" env:"HTTP_PORT"`
	}
)

func NewConfig() *Config {
	cfg := &Config{}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatalf("config path is not set")
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalf("failed to get absolute path: %s, path: %s", err, configPath)
	}
	err = cleanenv.ReadConfig(absConfigPath, cfg)
	if err != nil {
		log.Fatalf("failed to get absolute path: %s, path: %s", err, configPath)
	}
	err = cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("config error: %s", err)
	}

	return cfg
}

func (cfg *Config) GetInstance() *Config {
	cfg.once.Do(func() {
		cfg.instance = NewConfig()
	})

	return cfg.instance
}

func (cfg *Config) String() string {
	return fmt.Sprintf("App: %s, Version: %s, GRPC Port: %d, HTTP Port: %d",
		cfg.App.Name,
		cfg.App.Version,
		cfg.Grpc.Port,
		cfg.Http.Port,
	)
}

var ConfigSingleton Config
