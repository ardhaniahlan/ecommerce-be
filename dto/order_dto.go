package dto

type AdminOrderItem struct {
	ProductID   string  `db:"product_id" json:"product_id"`
	ProductName string  `db:"product_name" json:"product_name"`
	UnitPrice   float64 `db:"unit_price" json:"unit_price"`
	Quantity    int     `db:"quantity" json:"quantity"`
}

type AdminOrderResponse struct {
	ID              string           `db:"id" json:"id"`
	UserID          string           `db:"user_id" json:"user_id"`
	GrossAmount     float64          `db:"gross_amount" json:"gross_amount"`
	PaymentStatus   string           `db:"payment_status" json:"payment_status"`
	ShippingAddress *string          `db:"shipping_address" json:"shipping_address"`
	TrackingNumber  *string          `db:"tracking_number" json:"tracking_number"`
	
	Items           []AdminOrderItem `json:"items"` 
}

type TrackingRequest struct {
	TrackingNumber string `json:"tracking_number" binding:"required"`
}

type AdminDashboardStats struct {
	TotalRevenue    float64 `json:"total_revenue"`
	OrdersToProcess int     `json:"orders_to_process"`
	TotalProducts   int     `json:"total_products"`
	TotalUsers      int     `json:"total_users"`
}