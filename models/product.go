package models

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}
