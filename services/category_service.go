package services

import (
	"github.com/pisnov/golang_kasir/models"
	"github.com/pisnov/golang_kasir/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(c *models.Category) error {
	return s.repo.Create(c)
}

func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

func (s *CategoryService) Update(c *models.Category) error {
	return s.repo.Update(c)
}

func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}
