package services

import (
	"database/sql"
	"log"

	database "github.com/adewoleadenigbagbe/simpleloadbalancer/server/db"
)

type GetProductByIdRequest struct {
	Id string `param:"Id"`
}
type GetProductByIdResponse struct {
	Product Product `json:"product"`
}

func (productService ProductService) GetProductById(request GetProductByIdRequest) GetProductByIdResponse {
	row := database.DB.QueryRow("SELECT * FROM products WHERE Id=?", request.Id)

	var product Product
	var err error
	if err = row.Scan(&product.Id, &product.Price, &product.UnitInStock, &product.Name, &product.Discontinued, &product.Category, &product.CreatedOn, &product.ModifiedOn); err == sql.ErrNoRows {
		log.Fatal(err)
	}

	return GetProductByIdResponse{Product: product}
}
