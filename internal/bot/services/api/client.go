package api

import (
	"errors"
	"net/http"
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
	return "", errors.New("метод регистрации еще не реализован")
}
