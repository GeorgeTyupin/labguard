package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Server HTTPServerConf `yaml:"http_server"`
}

type HTTPServerConf struct {
	Address   string       `yaml:"address" env-default:"localhost:8080"`
	JWTSecret string       `env:"JWT_SECRET" env-required:"true"`
	Timeouts  TimeoutsConf `yaml:"timeouts"`
}

type TimeoutsConf struct {
	Idle     time.Duration `yaml:"idle" env-default:"60s"`
	Request  time.Duration `yaml:"request" env-default:"5s"`
	Shutdown time.Duration `yaml:"shutdown" env-default:"10s"`
}

func LoadServerConf(file *os.File) (*ServerConfig, error) {
	var config ServerConfig

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("не удалось прочитать env. Возникла ошибка %w", err)
	}

	if err := cleanenv.ParseYAML(file, &config); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("не удалось прочитать env переменные. Возникла ошибка %w", err)
	}

	return &config, nil
}
