package dto

type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartInput struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}