package dto

import "time"

type CreateProductInput struct {
	Name               string     `json:"name" binding:"required"`
	Description        *string    `json:"description"`
	Price              float64    `json:"price" binding:"required,min=0"`
	Stock              int        `json:"stock" binding:"required,min=0"`
	ImageURLs          []string   `json:"image_url"`
	DiscountPercentage int        `json:"discount_percentage"`
	DiscountStart      *time.Time `json:"discount_start"`
	DiscountEnd        *time.Time `json:"discount_end"`
}
