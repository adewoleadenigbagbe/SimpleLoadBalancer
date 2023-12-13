package services

type GetProductByIdRequest struct {
	Id string `param:"Id"`
}
type GetProductByIdResponse struct {
	Product Product `json:"product"`
}

func (productService ProductService) GetProductById(request GetProductByIdRequest) GetProductByIdResponse {
	return GetProductByIdResponse{}
}
