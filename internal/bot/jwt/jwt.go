package jwt

import (
	"fmt"
	"time"

	"github.com/GeorgeTyupin/labguard/internal/bot/config"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(cfg *config.Config) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["bot_name"] = cfg.BotName
	claims["exp"] = time.Now().Add(cfg.Client.JWT.TokenTTL).Unix()

	// Подписываем токен секретом
	tokenString, err := token.SignedString([]byte(cfg.Client.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи jwt: %w", err)
	}

	return tokenString, nil
}
