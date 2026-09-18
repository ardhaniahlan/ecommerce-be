package repositories

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type CartItemResponse struct {
	ID          int     `db:"id" json:"id"`
	ProductID   int     `db:"product_id" json:"productId"`
	ProductName string  `db:"name" json:"productName"`
	Price       float64 `db:"price" json:"price"`
	ImageURL    *string `db:"image_url" json:"imageUrl"`
	Quantity    int     `db:"quantity" json:"quantity"`
}

type CartRepository interface {
	AddToCart(userID string, productID int, quantity int) error
	GetCartByUserID(userID string) ([]CartItemResponse, error)
	RemoveFromCart(cartID int, userID string) error
	UpdateCartItemQuantity(cartID int, userID string, quantity int) error
}

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) CartRepository {
	return &cartRepository{db}
}

func (r *cartRepository) AddToCart(userID string, productID int, quantity int) error {
	var isActive bool
	err := r.db.Get(&isActive, `SELECT is_active FROM products WHERE id = $1`, productID)

	if err == sql.ErrNoRows {
		return errors.New("produk tidak ditemukan")
	} else if err != nil {
		return err
	}

	if !isActive {
		return errors.New("produk sudah tidak tersedia atau telah dihapus")
	}

	var existingID, existingQty int

	checkQuery := `SELECT id, quantity FROM cart_items WHERE user_id = $1 AND product_id = $2`
	err = r.db.QueryRow(checkQuery, userID, productID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)`
		_, err = r.db.Exec(insertQuery, userID, productID, quantity)
		return err
	} else if err != nil {
		return err
	}

	newQty := existingQty + quantity
	updateQuery := `UPDATE cart_items SET quantity = $1 WHERE id = $2`
	_, err = r.db.Exec(updateQuery, newQty, existingID)
	return err
}

func (r *cartRepository) GetCartByUserID(userID string) ([]CartItemResponse, error) {
	var items []CartItemResponse

	query := `
		SELECT c.id, c.product_id, p.name, p.price, p.image_url, c.quantity 
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
		ORDER BY c.id DESC`

	err := r.db.Select(&items, query, userID)
	return items, err
}

func (r *cartRepository) RemoveFromCart(cartID int, userID string) error {
	query := `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, cartID, userID)
	return err
}

func (r *cartRepository) UpdateCartItemQuantity(cartID int, userID string, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1 WHERE id = $2 AND user_id = $3`

	result, err := r.db.Exec(query, quantity, cartID, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("barang tidak ditemukan di keranjang Anda")
	}

	return nil
}
