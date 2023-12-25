package handlers

import (
	"net/http"

	"github.com/adewoleadenigbagbe/simpleloadbalancer/server/services"
	"github.com/labstack/echo/v4"
)

func AddProductHandler(productContext echo.Context) error {
	req := new(services.Product)
	err := productContext.Bind(req)
	if err != nil {
		return productContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	productService := services.ProductService{}
	resp := productService.AddProduct(*req)

	if resp == "" {
		return productContext.JSON(http.StatusInternalServerError, "could not save product")
	}
	return productContext.JSON(http.StatusOK, resp)
}

func GetProductsHandler(productContext echo.Context) error {
	req := new(services.GetProductRequest)
	err := productContext.Bind(req)
	if err != nil {
		return productContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	productService := services.ProductService{}
	resp := productService.GetProducts(*req)

	return productContext.JSON(http.StatusOK, resp)
}

func GetProductByIdHandler(productContext echo.Context) error {
	req := new(services.GetProductByIdRequest)
	err := productContext.Bind(req)
	if err != nil {
		return productContext.JSON(http.StatusBadRequest, "Bad Request")
	}

	productService := services.ProductService{}
	resp, err := productService.GetProductById(*req)
	if err != nil {
		return productContext.JSON(http.StatusInternalServerError, err.Error())
	}

	return productContext.JSON(http.StatusOK, resp)
}
