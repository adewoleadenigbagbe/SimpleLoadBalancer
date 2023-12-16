package services

import (
	"database/sql"
	"fmt"
	"log"

	database "github.com/adewoleadenigbagbe/simpleloadbalancer/server/db"
)

type GetProductRequest struct {
	Page       int    `query:"page"`
	PageLength int    `query:"pageLength"`
	SortBy     string `query:"sortBy"`
	Order      string `query:"order"`
}

type GetProductResponse struct {
	Products []Product
}

func (productService ProductService) GetProducts(request GetProductRequest) GetProductResponse {
	if request.Page < 0 {
		request.Page = 0
	}

	if request.PageLength < 1 {
		request.PageLength = 10
	}

	offset := request.Page * request.PageLength
	orderby := fmt.Sprint(request.SortBy, " ", request.Order)
	rows, err := database.DB.Query("SELECT * FROM products ORDER BY ? LIMIT ? OFFSET ?", orderby, request.PageLength, offset)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err = rows.Scan(&product.Id, &product.Price, &product.UnitInStock, &product.Name, &product.Discontinued, &product.Category, &product.CreatedOn, &product.ModifiedOn); err == sql.ErrNoRows {
			log.Fatal(err)
		}
		products = append(products, product)
	}

	return GetProductResponse{Products: products}
}
