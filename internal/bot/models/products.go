package models

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	Purchased   bool
	Link        string
}
