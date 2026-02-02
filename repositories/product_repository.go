package repositories

import (
	"database/sql"
	"errors"

	"github.com/pisnov/golang_kasir/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id`
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		var catID sql.NullInt64
		var catName sql.NullString
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName)
		if err != nil {
			return nil, err
		}
		if catID.Valid {
			v := int(catID.Int64)
			p.CategoryID = &v
		} else {
			p.CategoryID = nil
		}
		if catName.Valid {
			p.CategoryName = catName.String
		} else {
			p.CategoryName = ""
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	return repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
}

func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, c.name
	FROM products p
	LEFT JOIN categories c ON p.category_id = c.id
	WHERE p.id = $1`
	var p models.Product
	var catID sql.NullInt64
	var catName sql.NullString
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &catID, &catName)
	if err == sql.ErrNoRows {
		return nil, errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	if catID.Valid {
		v := int(catID.Int64)
		p.CategoryID = &v
	}
	if catName.Valid {
		p.CategoryName = catName.String
	}
	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}
