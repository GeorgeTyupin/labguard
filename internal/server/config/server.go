package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Server HTTPServerConf `yaml:"http_server"`
}

type HTTPServerConf struct {
	Address   string       `yaml:"address" env-default:"localhost:8080"`
	JWTSecret string       `yaml:"jwt_secret" env-required:"true"`
	Timeouts  TimeoutsConf `yaml:"timeouts"`
}

type TimeoutsConf struct {
	Idle     time.Duration `yaml:"idle" env-default:"60s"`
	Request  time.Duration `yaml:"request" env-default:"5s"`
	Shutdown time.Duration `yaml:"shutdown" env-default:"10s"`
}

func LoadServerConf(file os.File) (*ServerConfig, error) {
	var config ServerConfig

	if err := cleanenv.ParseYAML(&file, &config); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	return &config, nil
}
