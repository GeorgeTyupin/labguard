package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresConfig struct {
	Postgres DBConfig `yaml:"postgres"`
}

type DBConfig struct {
	Database   string         `yaml:"database" env:"POSTGRES_DB"`
	User       string         `yaml:"user" env:"POSTGRES_USER"`
	Password   string         `yaml:"password" env:"POSTGRES_PASSWORD"`
	Host       string         `yaml:"host" env:"POSTGRES_HOST"`
	Port       int            `yaml:"port" env:"POSTGRES_PORT"`
	PoolSize   int32          `yaml:"pool_size" env-default:"10"`
	Connection ConnectionConf `yaml:"connection"`
}

type ConnectionConf struct {
	MaxLifeTime       time.Duration `yaml:"max_life_time" env-default:"30min"`
	MaxIdleTime       time.Duration `yaml:"max_idle_time" env-default:"1min"`
	HealthCheckPeriod time.Duration `yaml:"health_check_period" env-default:"30s"`
	Timeout           time.Duration `yaml:"timeout" env-default:"30s"`
}

func LoadPostgresConf(file os.File) (*PostgresConfig, error) {
	var pgConf PostgresConfig

	if err := cleanenv.ParseYAML(&file, &pgConf); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг. Возникла ошибка %w", err)
	}

	return &pgConf, nil

}
