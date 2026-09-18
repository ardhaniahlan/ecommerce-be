package models

import "time"

type User struct {
	ID            string    `db:"id" json:"id"`
	FullName      string    `db:"full_name" json:"fullName"`
	Email         string    `db:"email" json:"email"`
	PasswordHash  string    `db:"password_hash" json:"-"`
	Role          string    `db:"role" json:"role"`
	Phone         *string   `db:"phone" json:"phone,omitempty"`
	Province      *string   `db:"province" json:"province,omitempty"`
	City          *string   `db:"city" json:"city,omitempty"`
	District      *string   `db:"district" json:"district,omitempty"`
	PostalCode    *string   `db:"postal_code" json:"postalCode,omitempty"`
	StreetAddress *string   `db:"street_address" json:"streetAddress,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}

type Product struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	Price       float64   `db:"price" json:"price"`
	Stock       int       `db:"stock" json:"stock"`
	ImageURL    *string   `db:"image_url" json:"imageUrl,omitempty"`
	IsActive    bool      `db:"is_active" json:"isActive"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type CartItem struct {
	ID        int    `db:"id" json:"id"`
	UserID    string `db:"user_id" json:"userId"`
	ProductID int    `db:"product_id" json:"productId"`
	Quantity  int    `db:"quantity" json:"quantity"`
}

type Order struct {
	ID                    string    `db:"id" json:"id"`
	UserID                string    `db:"user_id" json:"userId"`
	OrderDate             time.Time `db:"order_date" json:"orderDate"`
	GrossAmount           float64   `db:"gross_amount" json:"grossAmount"`
	PaymentURL            *string   `db:"payment_url" json:"paymentUrl,omitempty"`
	PaymentStatus         string    `db:"payment_status" json:"paymentStatus"`
	MidtransTransactionID *string   `db:"midtrans_transaction_id" json:"midtransTransactionId,omitempty"`
	PaymentMethod         *string   `db:"payment_method" json:"paymentMethod,omitempty"`
}

type OrderItem struct {
	ID          int     `db:"id" json:"id"`
	OrderID     string  `db:"order_id" json:"orderId"`
	ProductID   int     `db:"product_id" json:"productId"`
	ProductName string  `db:"product_name" json:"productName"`
	UnitPrice   float64 `db:"unit_price" json:"unitPrice"`
	Quantity    int     `db:"quantity" json:"quantity"`
}
