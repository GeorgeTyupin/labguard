package models

type User struct {
	UUID     int64
	Name     string
	Group    string
	TokensID []int
}
