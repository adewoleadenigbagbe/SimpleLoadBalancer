package services

type GetProductRequest struct {
	Page       int
	PageLength int
	SortBy     string
	Order      string
}

type GetProductResponse struct {
	Products []Product
}

func (productService ProductService) GetProducts(request GetProductRequest) GetProductResponse {
	return GetProductResponse{}
}
