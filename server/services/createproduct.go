package services

import (
	"time"

	database "github.com/adewoleadenigbagbe/simpleloadbalancer/server/db"
	sequentialguid "github.com/adewoleadenigbagbe/simpleloadbalancer/server/helpers"
)

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
	today := time.Now()
	product.Id = sequentialguid.New().String()
	product.Discontinued = false
	product.CreatedOn = today
	product.ModifiedOn = today

	_, err := database.DB.Exec("INSERT INTO products VALUES(?,?,?,?,?,?,?,?);",
		product.Id, product.Price, product.UnitInStock, product.Name,
		product.Discontinued, product.Category, product.CreatedOn, product.ModifiedOn)

	if err != nil {
		return ""
	}
	return product.Id
}
