package services

import (
	"github.com/pisnov/golang_kasir/models"
	"github.com/pisnov/golang_kasir/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetSalesSummaryToday() (*models.SalesSummary, error) {
	return s.repo.GetSalesSummaryToday()
}

func (s *TransactionService) GetSalesSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error) {
	return s.repo.GetSalesSummaryByDateRange(startDate, endDate)
}
