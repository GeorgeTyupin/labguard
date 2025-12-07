package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Неверный формат header авторизации", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("неопределенный метод подписи: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			}, nil)

			if err != nil {
				http.Error(w, "Не удалось спарсить токен", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "невалидный токен", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
