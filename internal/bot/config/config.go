package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	yamlPath = "/configs/bot/bot.yaml"
	envPath  = "/configs/bot/bot.env"
)

type Config struct {
	BotConf `yaml:"bot"`
}

type BotConf struct {
	BotName  string        `yaml:"name"  env-default:"bot"`
	BotToken string        `env-required:"true" env:"BOT_TOKEN"`
	Client   BotClientConf `yaml:"client"`
}

type BotClientConf struct {
	ServerAddress string  `yaml:"server_address" env-default:"http://localhost:8080"`
	JWT           JWTConf `yaml:"jwt"`
}

type JWTConf struct {
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"30min"`
	Secret   string        `env:"JWT_SECRET"`
}

func MustLoad(logger *slog.Logger) *Config {
	const op = "bot.config.MustLoad"
	logger = logger.With(slog.String("op", op))
	cfg, err := newBotConf()
	if err != nil {
		logger.Error("Ошибка загрузки конфига бота", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return cfg
}

func newBotConf() (*Config, error) {
	file, err := os.Open(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("возникла ошибка с открытием файла конфига по пути %s, возникла ошибка %w", yamlPath, err)
	}
	defer file.Close()

	var cfg Config

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("не удалось прочитать env. Возникла ошибка %w", err)
	}

	if err := cleanenv.ParseYAML(file, &cfg); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось прочитать env переменные. Возникла ошибка %w", err)
	}

	return &cfg, nil
}
