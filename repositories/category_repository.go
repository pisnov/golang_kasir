package repositories

import (
	"database/sql"
	"errors"

	"github.com/pisnov/golang_kasir/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := repo.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

func (repo *CategoryRepository) Create(cat *models.Category) error {
	return repo.db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", cat.Name).Scan(&cat.ID)
}

func (repo *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := repo.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&c.ID, &c.Name)
	if err == sql.ErrNoRows {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *CategoryRepository) Update(cat *models.Category) error {
	result, err := repo.db.Exec("UPDATE categories SET name = $1 WHERE id = $2", cat.Name, cat.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}
	return nil
}

func (repo *CategoryRepository) Delete(id int) error {
	result, err := repo.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("kategori tidak ditemukan")
	}
	return nil
}
