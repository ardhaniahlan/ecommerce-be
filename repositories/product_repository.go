package repositories

import (
	"ecommerce-backend/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetAll() ([]models.Product, error)
	GetByID(id int) (models.Product, error)
	Update(product *models.Product) error
	Delete(id int) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, image_url) 
		VALUES (:name, :description, :price, :stock, :image_url) 
		RETURNING id, is_active, created_at`
	
	rows, err := r.db.NamedQuery(query, product)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(product)
	}
	return err
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	query := `SELECT id, name, description, price, stock, image_url, is_active, created_at FROM products WHERE is_active = true`
	
	err := r.db.Select(&products, query)
	return products, err
}

func (r *productRepository) GetByID(id int) (models.Product, error) {
	var product models.Product
	query := `SELECT id, name, description, price, stock, image_url, is_active, created_at FROM products WHERE id = $1`
	err := r.db.Get(&product, query, id)
	return product, err
}

func (r *productRepository) Update(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = :name, description = :description, price = :price, stock = :stock, image_url = :image_url 
		WHERE id = :id 
		RETURNING is_active, created_at`
	
	rows, err := r.db.NamedQuery(query, product)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(product)
	}
	return err
}

func (r *productRepository) Delete(id int) error {
	query := `UPDATE products SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}