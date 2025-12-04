package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const confPath = "/configs/server/server.yaml"

type Config struct {
	Env    string         `yaml:"env" env-default:"local"`
	Server HTTPServerConf `yaml:"http_server"`
}

type HTTPServerConf struct {
	Address  string       `yaml:"address" env-default:"localhost:8080"`
	Timeouts TimeoutsConf `yaml:"timeouts"`
}

type TimeoutsConf struct {
	Idle     time.Duration `yaml:"idle" env-default:"60s"`
	Request  time.Duration `yaml:"request" env-default:"5s"`
	Shutdown time.Duration `yaml:"shutdown" env-default:"10s"`
}

func LoadConf() (*Config, error) {
	// TODO Сделать загрузку пути к конфигу из переменных окружения
	var config Config

	file, err := os.Open(confPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл с конфигом. Возникла ошибка %w", err)
	}

	err = cleanenv.ParseYAML(file, &config)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	return &config, nil
}
