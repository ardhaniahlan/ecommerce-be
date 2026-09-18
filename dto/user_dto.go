package dto

type UpdateProfileRequest struct {
	FullName      string `json:"full_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	Province      string `json:"province" binding:"required"`
	City          string `json:"city" binding:"required"`
	District      string `json:"district" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	StreetAddress string `json:"street_address" binding:"required"`
}
