package repositories

import (
	"ecommerce-backend/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(product *models.Product, imageURLs []string) error
	GetAll(page int, limit int, search string) ([]models.Product, int, error)
	GetByID(id int) (models.Product, error)
	Update(product *models.Product, imageURLs []string) error
	Delete(id int) error
}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *models.Product, imageURLs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	queryProduct := `
		INSERT INTO products (name, description, price, stock, discount_percentage, discount_start, discount_end) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id
	`

	err = tx.QueryRow(queryProduct,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.DiscountPercentage,
		product.DiscountStart,
		product.DiscountEnd,
	).Scan(&product.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	queryImage := `INSERT INTO product_images (product_id, image_url, is_primary) VALUES ($1, $2, $3)`

	for i, url := range imageURLs {
		isPrimary := false
		if i == 0 {
			isPrimary = true
		}

		if _, err = tx.Exec(queryImage, product.ID, url, isPrimary); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *productRepository) GetAll(page int, limit int, search string) ([]models.Product, int, error) {
	var products []models.Product
	var totalItems int

	offset := (page - 1) * limit

	searchParam := "%" + search + "%"

	countQuery := `SELECT count(*) FROM products WHERE is_active = true AND name ILIKE $1`
	err := r.db.Get(&totalItems, countQuery, searchParam)
	if err != nil {
		return nil, 0, err
	}

	query := `
  SELECT p.id, p.name, p.description, p.price, p.stock, p.is_active, p.created_at,
         p.discount_percentage, p.discount_start, p.discount_end,
         pi.image_url 
  FROM products p
  LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
  WHERE p.is_active = true AND p.name ILIKE $1 
  ORDER BY p.created_at DESC
  LIMIT $2 OFFSET $3`

	err = r.db.Select(&products, query, searchParam, limit, offset)
	return products, totalItems, err
}

func (r *productRepository) GetByID(id int) (models.Product, error) {
	var product models.Product

	queryProduct := `SELECT id, name, description, price, stock, discount_percentage, discount_start, discount_end, is_active, created_at FROM products WHERE id = $1`
	err := r.db.Get(&product, queryProduct, id)
	if err != nil {
		return product, err
	}

	var images []models.ProductImage
	queryImages := `SELECT id, product_id, image_url, is_primary FROM product_images WHERE product_id = $1 ORDER BY is_primary DESC`

	err = r.db.Select(&images, queryImages, id)
	if err != nil {
		return product, err
	}

	product.Images = images

	return product, nil
}

func (r *productRepository) Update(product *models.Product, imageURLs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	queryProduct := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, stock = $4, 
		    discount_percentage = $5, discount_start = $6, discount_end = $7
		WHERE id = $8 
		RETURNING is_active, created_at`

	err = tx.QueryRow(queryProduct,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.DiscountPercentage,
		product.DiscountStart,
		product.DiscountEnd,
		product.ID,
	).Scan(&product.IsActive, &product.CreatedAt)

	if err != nil {
		tx.Rollback()
		return err
	}

	if len(imageURLs) > 0 {
		queryDeleteImages := `DELETE FROM product_images WHERE product_id = $1`
		if _, err = tx.Exec(queryDeleteImages, product.ID); err != nil {
			tx.Rollback()
			return err
		}

		queryInsertImage := `INSERT INTO product_images (product_id, image_url, is_primary) VALUES ($1, $2, $3)`
		for i, url := range imageURLs {
			isPrimary := false
			if i == 0 {
				isPrimary = true
			}
			if _, err = tx.Exec(queryInsertImage, product.ID, url, isPrimary); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *productRepository) Delete(id int) error {
	query := `UPDATE products SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func calculateActivePrice(originalPrice float64, discountPercentage int, start, end *time.Time) (activePrice float64, validPercentage int) {
	if discountPercentage > 0 && start != nil && end != nil {
		now := time.Now()
		if now.After(*start) && now.Before(*end) {
			discountAmount := (originalPrice * float64(discountPercentage)) / 100
			return originalPrice - discountAmount, discountPercentage
		}
	}
	return originalPrice, 0
}
