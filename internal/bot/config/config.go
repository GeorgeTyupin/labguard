package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// TODO: добавить настройки когда понадобятся
}

func Load(confPath string) (*Config, error) {
	// TODO: раскомментировать когда появятся настройки в YAML
	// file, err := os.Open(confPath)
	// if err != nil {
	// 	return nil, fmt.Errorf("возникла ошибка с открытием файла конфига по пути %s, возникла ошибка %w", confPath, err)
	// }
	// defer file.Close()
	//
	// var cfg Config
	//
	// decoder := yaml.NewDecoder(file)
	// if err := decoder.Decode(&cfg); err != nil {
	// 	return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	// }
	//
	// return &cfg, nil

	// Пока возвращаем пустую структуру
	return &Config{}, nil
}

func GetBotToken(envPath string) (string, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return "", fmt.Errorf("не удалось загрузить переменные среды по пути %s, возникла ошибка %w", envPath, err)
	}

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return "", fmt.Errorf("переменная окружения BOT_TOKEN не установлена")
	}

	return token, nil
}
