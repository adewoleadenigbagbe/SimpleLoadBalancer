package services

import "time"

type Product struct {
	Id           string
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	UnitInStock  int     `json:"unit"`
	Discontinued bool    `json:"discontinued"`
	CreatedOn    time.Time
	ModifiedOn   time.Time
}

func (productService ProductService) AddProduct(product Product) string {
	return ""
}
