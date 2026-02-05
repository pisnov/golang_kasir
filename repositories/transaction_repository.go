package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pisnov/golang_kasir/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Check stock
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", productName)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Update product stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetSalesSummaryToday() (*models.SalesSummary, error) {
	today := time.Now().Format("2006-01-02")
	return repo.GetSalesSummaryByDateRange(today, today)
}

func (repo *TransactionRepository) GetSalesSummaryByDateRange(startDate, endDate string) (*models.SalesSummary, error) {
	// Total revenue and transaction count
	var totalRevenue int
	var totalTransaksi int

	query := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) >= $1 AND DATE(created_at) <= $2
	`
	err := repo.db.QueryRow(query, startDate, endDate).Scan(&totalRevenue, &totalTransaksi)
	if err != nil {
		return nil, err
	}

	// Best selling product
	var bestProduct *models.BestSellingItem

	queryBest := `
		SELECT p.name, SUM(td.quantity)
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) >= $1 AND DATE(t.created_at) <= $2
		GROUP BY p.name
		ORDER BY SUM(td.quantity) DESC
		LIMIT 1
	`
	var name string
	var qty int
	err = repo.db.QueryRow(queryBest, startDate, endDate).Scan(&name, &qty)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err != sql.ErrNoRows {
		bestProduct = &models.BestSellingItem{
			Nama:       name,
			QtyTerjual: qty,
		}
	}

	return &models.SalesSummary{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: bestProduct,
	}, nil
}
