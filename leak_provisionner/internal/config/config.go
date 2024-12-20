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
		Name       string `yaml:"name"`
		Version    string `yaml:"version"`
		Log        `yaml:"log"`
		GrpcServer `yaml:"grpc_server"`
		FileReader `yaml:"file_reader"`
	}

	FileReader struct {
		PasswordBatchSize int `yaml:"password_batch_size" env:"FILE_READER_PASSWORD_BATCH_SIZE"`
	}

	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	}

	GrpcServer struct {
		Host string `yaml:"host" env:"GRPC_SERVER_HOST"`
		Port int    `yaml:"port" env:"GRPC_SERVER_PORT"`
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
	return fmt.Sprintf("App: %s, Version: %s, GRPC Host: %s, GRPC Port: %d, Password Batch Size: %d",
		cfg.App.Name,
		cfg.App.Version,
		cfg.GrpcServer.Host,
		cfg.GrpcServer.Port,
		cfg.FileReader.PasswordBatchSize,
	)
}

var ConfigSingleton Config
