package dto

type AdminOrderResponse struct {
	ID              string  `db:"id" json:"id"`
	UserID          string  `db:"user_id" json:"user_id"`
	GrossAmount     float64 `db:"gross_amount" json:"gross_amount"`
	PaymentStatus   string  `db:"payment_status" json:"payment_status"`
	ShippingAddress *string `db:"shipping_address" json:"shipping_address"`
	TrackingNumber  *string `db:"tracking_number" json:"tracking_number"`
}

type TrackingRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
}