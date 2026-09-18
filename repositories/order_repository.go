package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Checkout(userID string) (string, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Checkout(userID string) (orderID string, err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return "", err
	}

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
		return "", err
	}

	if len(cartItems) == 0 {
		return "", errors.New("keranjang belanja kosong")
	}

	var grossAmount float64
	for _, item := range cartItems {
		if item.Quantity > item.Stock {
			return "", fmt.Errorf("stok produk '%s' tidak mencukupi (sisa: %d, diminta: %d)", item.ProductName, item.Stock, item.Quantity)
		}
		grossAmount += (item.UnitPrice * float64(item.Quantity))
	}

	orderID = fmt.Sprintf("ORD-%d", time.Now().Unix())

	queryOrder := `
		INSERT INTO orders (id, user_id, gross_amount, payment_status)
		VALUES ($1, $2, $3, 'Unpaid')
	`
	if _, err = tx.Exec(queryOrder, orderID, userID, grossAmount); err != nil {
		return "", err
	}

	queryOrderItem := `INSERT INTO order_items (order_id, product_id, product_name, unit_price, quantity) VALUES ($1, $2, $3, $4, $5)`
	queryUpdateStock := `UPDATE products SET stock = stock - $1 WHERE id = $2`
	queryDeleteCart := `DELETE FROM cart_items WHERE id = $1`

	for _, item := range cartItems {
		if _, err = tx.Exec(queryOrderItem, orderID, item.ProductID, item.ProductName, item.UnitPrice, item.Quantity); err != nil {
			return "", err
		}

		if _, err = tx.Exec(queryUpdateStock, item.Quantity, item.ProductID); err != nil {
			return "", err
		}

		if _, err = tx.Exec(queryDeleteCart, item.CartID); err != nil {
			return "", err
		}
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return orderID, nil
}
