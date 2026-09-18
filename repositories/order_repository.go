package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Checkout(userID string) (string, float64, error)
	UpdatePaymentURL(orderID string, paymentURL string) error
	UpdateOrderStatus(orderID string, status string, midtransID, payMethod string) error
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Checkout(userID string) (orderID string, grossAmount float64, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return "", 0, err
	}

	var user struct {
		Phone         *string `db:"phone"`
		Province      *string `db:"province"`
		City          *string `db:"city"`
		District      *string `db:"district"`
		PostalCode    *string `db:"postal_code"`
		StreetAddress *string `db:"street_address"`
	}

	err = tx.Get(&user, "SELECT phone, province, city, district, postal_code, street_address FROM users WHERE id = $1", userID)
	if err != nil {
		return "", 0, err
	}

	if user.Phone == nil || user.Province == nil || user.StreetAddress == nil {
		return "", 0, errors.New("alamat pengiriman belum lengkap, harap update profil Anda terlebih dahulu")
	}

	fullShippingAddress := fmt.Sprintf("%s, %s, %s, %s, %s. (HP: %s)",
		*user.StreetAddress, *user.District, *user.City, *user.Province, *user.PostalCode, *user.Phone)

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var cartItems []struct {
		CartID      int     `db:"cart_id"`
		ProductID   int     `db:"product_id"`
		ProductName string  `db:"name"`
		UnitPrice   float64 `db:"price"`
		Quantity    int     `db:"quantity"`
		Stock       int     `db:"stock"`
	}

	queryCart := `
		SELECT c.id as cart_id, p.id as product_id, p.name, p.price, c.quantity, p.stock
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
	`
	if err = tx.Select(&cartItems, queryCart, userID); err != nil {
		return "", 0, err
	}

	if len(cartItems) == 0 {
		return "", 0, errors.New("keranjang belanja kosong")
	}

	for _, item := range cartItems {
		if item.Quantity > item.Stock {
			return "", 0, fmt.Errorf("stok produk '%s' tidak mencukupi", item.ProductName)
		}
		grossAmount += (item.UnitPrice * float64(item.Quantity))
	}

	orderID = fmt.Sprintf("ORD-%d", time.Now().Unix())

	queryOrder := `
		INSERT INTO orders (id, user_id, gross_amount, payment_status, shipping_address) 
		VALUES ($1, $2, $3, 'Unpaid', $4)
	`
	if _, err = tx.Exec(queryOrder, orderID, userID, grossAmount, fullShippingAddress); err != nil {
		return "", 0, err
	}

	queryOrderItem := `INSERT INTO order_items (order_id, product_id, product_name, unit_price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryDeleteCart := `DELETE FROM cart_items WHERE id = $1`
	queryUpdateStock := `UPDATE products SET stock = stock - $1 WHERE id = $2 AND stock >= $1`

	for _, item := range cartItems {
		if _, err = tx.Exec(queryOrderItem, orderID, item.ProductID, item.ProductName, item.UnitPrice, item.Quantity); err != nil {
			return "", 0, err
		}

		result, err := tx.Exec(queryUpdateStock, item.Quantity, item.ProductID)
		if err != nil {
			return "", 0, err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return "", 0, fmt.Errorf("checkout dibatalkan: ada user lain yang baru saja memborong '%s', stok tidak lagi mencukupi", item.ProductName)
		}

		if _, err = tx.Exec(queryDeleteCart, item.CartID); err != nil {
			return "", 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return "", 0, err
	}

	return orderID, grossAmount, nil
}

func (r *orderRepository) UpdatePaymentURL(orderID string, paymentURL string) error {
	query := `UPDATE orders SET payment_url = $1 WHERE id = $2`
	_, err := r.db.Exec(query, paymentURL, orderID)
	return err
}

func (r *orderRepository) UpdateOrderStatus(orderID string, status string, midtransID, payMethod string) error {
	query := `UPDATE orders 
	          SET payment_status = $1, midtrans_transaction_id = $2, payment_method = $3 
	          WHERE id = $4`

	result, err := r.db.Exec(query, status, midtransID, payMethod, orderID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("pesanan dengan ID %s tidak ditemukan di database", orderID)
	}

	return nil
}
