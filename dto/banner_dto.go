package dto

type BannerInput struct {
	Title    string  `json:"title" binding:"required"`
	ImageURL string  `json:"imageUrl" binding:"required"`
}