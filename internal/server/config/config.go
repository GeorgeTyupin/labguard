package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	yamlPath = "configs/server/server.yaml"
	envPath  = "configs/server/postgres.env"
)

type Config struct {
	Env string
	ServerConfig
	PostgresConfig
}

func MustLoad(logger *slog.Logger) *Config {
	const op = "server.config.MustLoad"
	logger = logger.With(slog.String("op", op))

	file, err := os.Open(yamlPath)
	if err != nil {
		logger.Error("Не удалось открыть файл с конфигом", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer file.Close()

	serverConf, err := LoadServerConf(file)
	if err != nil {
		logger.Error("Ошибка загрузки конфига сервера", slog.String("error", err.Error()))
		os.Exit(1)
	}

	file.Seek(0, 0)
	postgresConf, err := LoadPostgresConf(file)
	if err != nil {
		logger.Error("Ошибка загрузки конфига базы данных", slog.String("error", err.Error()))
		os.Exit(1)
	}

	file.Seek(0, 0)
	envConf, err := LoadEnvState(file)
	if err != nil {
		logger.Error("Ошибка загрузки конфига env", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return &Config{
		Env:            envConf,
		ServerConfig:   *serverConf,
		PostgresConfig: *postgresConf,
	}
}

func LoadEnvState(file *os.File) (string, error) {
	var cfg struct {
		Env string `yaml:"env" env-default:"local"`
	}

	if err := cleanenv.ParseYAML(file, &cfg); err != nil {
		return "", fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	return cfg.Env, nil
}
