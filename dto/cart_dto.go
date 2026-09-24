package dto

import "time"

type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type UpdateCartInput struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

type CartItemResponse struct {
	ID                 int        `db:"id" json:"id"`
	ProductID          int        `db:"product_id" json:"productId"`
	ProductName        string     `db:"name" json:"productName"`
	OriginalPrice      float64    `db:"price" json:"originalPrice"`
	ActivePrice        float64    `json:"activePrice"`
	DiscountPercentage int        `db:"discount_percentage" json:"discountPercentage"`
	DiscountStart      *time.Time `db:"discount_start" json:"-"`
	DiscountEnd        *time.Time `db:"discount_end" json:"-"`
	ImageURL           *string    `db:"image_url" json:"imageUrl"`
	Quantity           int        `db:"quantity" json:"quantity"`
	Subtotal           float64    `json:"subtotal"`
	Stock              int        `json:"stock"`
}