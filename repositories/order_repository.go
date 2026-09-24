package repositories

import (
	"database/sql"
	"ecommerce-backend/dto"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Checkout(userID string, voucherCode string) (string, float64, error)
	UpdatePaymentURL(orderID string, paymentURL string) error
	UpdateOrderStatus(orderID string, status string, midtransID, payMethod string) error

	GetAllOrders() ([]dto.AdminOrderResponse, error)
	UpdateTrackingNumber(orderID, trackingNumber string) error

	GetUserOrders(userID string) ([]dto.AdminOrderResponse, error)
	CompleteOrder(orderID, userID string) error

	GetDashboardStats() (dto.AdminDashboardStats, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Checkout(userID string, voucherCode string) (orderID string, grossAmount float64, err error) {
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
		CartID             int        `db:"cart_id"`
		ProductID          int        `db:"product_id"`
		ProductName        string     `db:"name"`
		OriginalPrice      float64    `db:"price"`
		Quantity           int        `db:"quantity"`
		Stock              int        `db:"stock"`
		DiscountPercentage int        `db:"discount_percentage"`
		DiscountStart      *time.Time `db:"discount_start"`
		DiscountEnd        *time.Time `db:"discount_end"`
		ActivePrice        float64
	}

	queryCart := `
		SELECT c.id as cart_id, p.id as product_id, p.name, p.price, c.quantity, p.stock,
		       p.discount_percentage, p.discount_start, p.discount_end
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

	for i := range cartItems {
		if cartItems[i].Quantity > cartItems[i].Stock {
			err = fmt.Errorf("stok produk '%s' tidak mencukupi", cartItems[i].ProductName)
			return "", 0, err
		}

		activePrice, _ := calculateActivePrice(
			cartItems[i].OriginalPrice,
			cartItems[i].DiscountPercentage,
			cartItems[i].DiscountStart,
			cartItems[i].DiscountEnd,
		)

		cartItems[i].ActivePrice = activePrice
		grossAmount += (activePrice * float64(cartItems[i].Quantity))
	}

	if voucherCode != "" {
		now := time.Now()

		var voucher struct {
			ID                 string     `db:"id"`
			DiscountPercentage int        `db:"discount_percentage"`
			MaxDiscountAmount  float64    `db:"max_discount_amount"`
			Quota              int        `db:"quota"`
			ValidFrom          *time.Time `db:"valid_from"`
			ValidUntil         *time.Time `db:"valid_until"`
		}

		err = tx.Get(&voucher, "SELECT id, discount_percentage, max_discount_amount, quota, valid_from, valid_until FROM vouchers WHERE code = $1", voucherCode)
		if err != nil {
			err = errors.New("kode voucher tidak valid atau tidak ditemukan")
			return "", 0, err
		}

		if voucher.Quota <= 0 {
			err = errors.New("kuota voucher sudah habis")
			return "", 0, err
		}

		if voucher.ValidFrom != nil && voucher.ValidUntil != nil {
			if now.Before(*voucher.ValidFrom) || now.After(*voucher.ValidUntil) {
				err = errors.New("masa berlaku voucher sudah habis atau belum dimulai")
				return "", 0, err
			}
		}

		discountApplied := (grossAmount * float64(voucher.DiscountPercentage)) / 100

		if voucher.MaxDiscountAmount > 0 && discountApplied > voucher.MaxDiscountAmount {
			discountApplied = voucher.MaxDiscountAmount
		}

		grossAmount -= discountApplied

		result, err := tx.Exec("UPDATE vouchers SET quota = quota - 1 WHERE id = $1 AND quota > 0", voucher.ID)
		if err != nil {
			return "", 0, err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			err = errors.New("checkout dibatalkan: kuota voucher baru saja habis digunakan orang lain")
			return "", 0, err
		}
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
		if _, err = tx.Exec(queryOrderItem, orderID, item.ProductID, item.ProductName, item.ActivePrice, item.Quantity); err != nil {
			return "", 0, err
		}

		var result sql.Result
		result, err = tx.Exec(queryUpdateStock, item.Quantity, item.ProductID)
		if err != nil {
			return "", 0, err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			err = fmt.Errorf("checkout dibatalkan: ada user lain yang baru saja memborong '%s', stok tidak lagi mencukupi", item.ProductName)
			return "", 0, err
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

func (r *orderRepository) GetAllOrders() ([]dto.AdminOrderResponse, error) {
	var orders []dto.AdminOrderResponse
	queryOrders := `
		SELECT id, user_id, gross_amount, payment_status, shipping_address, tracking_number 
		FROM orders ORDER BY order_date DESC
	`
	err := r.db.Select(&orders, queryOrders)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		var items []dto.AdminOrderItem
		queryItems := `
			SELECT product_id, product_name, unit_price, quantity 
			FROM order_items WHERE order_id = $1
		`
		err = r.db.Select(&items, queryItems, orders[i].ID)
		if err == nil {
			orders[i].Items = items
		} else {
			orders[i].Items = []dto.AdminOrderItem{}
		}
	}

	return orders, err
}

func (r *orderRepository) UpdateTrackingNumber(orderID, trackingNumber string) error {
	query := `
		UPDATE orders 
		SET tracking_number = $1, payment_status = 'Shipped' 
		WHERE id = $2 AND payment_status = 'Paid'
	`
	result, err := r.db.Exec(query, trackingNumber, orderID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("pesanan tidak ditemukan atau belum dibayar (status bukan Paid)")
	}

	return nil
}

func (r *orderRepository) GetUserOrders(userID string) ([]dto.AdminOrderResponse, error) {
	var orders []dto.AdminOrderResponse

	queryOrders := `
		SELECT id, user_id, gross_amount, payment_status, shipping_address, tracking_number 
		FROM orders 
		WHERE user_id = $1 
		ORDER BY id DESC
	`
	err := r.db.Select(&orders, queryOrders, userID)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		var items []dto.AdminOrderItem
		queryItems := `
			SELECT product_id, product_name, unit_price, quantity 
			FROM order_items WHERE order_id = $1
		`
		err = r.db.Select(&items, queryItems, orders[i].ID)
		if err == nil {
			orders[i].Items = items
		} else {
			orders[i].Items = []dto.AdminOrderItem{}
		}
	}

	return orders, nil
}

func (r *orderRepository) CompleteOrder(orderID, userID string) error {
	query := `
		UPDATE orders 
		SET payment_status = 'Completed' 
		WHERE id = $1 AND user_id = $2 AND payment_status = 'Shipped'
	`
	result, err := r.db.Exec(query, orderID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("pesanan tidak ditemukan atau barang belum dikirim (status bukan Shipped)")
	}

	return nil
}

func (r *orderRepository) GetDashboardStats() (dto.AdminDashboardStats, error) {
	var stats dto.AdminDashboardStats

	err := r.db.Get(&stats.TotalRevenue, `
		SELECT COALESCE(SUM(gross_amount), 0) 
		FROM orders 
		WHERE payment_status IN ('Paid', 'Shipped', 'Completed')
	`)
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&stats.OrdersToProcess, "SELECT COUNT(id) FROM orders WHERE payment_status = 'Paid'")
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&stats.TotalProducts, "SELECT COUNT(id) FROM products")
	if err != nil {
		return stats, err
	}

	err = r.db.Get(&stats.TotalUsers, "SELECT COUNT(id) FROM users WHERE role = 'user'")
	if err != nil {
		return stats, err
	}

	return stats, nil
}
