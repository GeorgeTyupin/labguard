package api

import (
	"net/http"

	"github.com/GeorgeTyupin/labguard/internal/bot/models"
)

// TODO Реализовать http клиент

type HttpClient struct {
	Client *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (client *HttpClient) CheckUserExists(uuid int64) (bool, error) {
	// TODO Реализовать этот метод после написания сервера
	return false, nil
}

func (client *HttpClient) RegisterUser(uuid int64, name, group string) (string, error) {
	// TODO Реализовать этот метод после написания сервера
	return "", nil
}

func (client *HttpClient) GetProducts(uuid int64) ([]*models.Product, error) {
	// TODO Реализовать реальный запрос для получения списка продуктов

	// Мок списка продуктов
	return []*models.Product{
		{
			ID:        1,
			Name:      "Лабораторная работа №1",
			Price:     500.00,
			Purchased: false,
		},
		{
			ID:        2,
			Name:      "Лабораторная работа №2",
			Price:     600.00,
			Purchased: true,
		},
		{
			ID:        3,
			Name:      "Курсовая работа",
			Price:     1500.00,
			Purchased: false,
		},
	}, nil
}
