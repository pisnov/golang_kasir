package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pisnov/golang_kasir/services"
)

type ReportHandler struct {
	transactionService *services.TransactionService
}

func NewReportHandler(transactionService *services.TransactionService) *ReportHandler {
	return &ReportHandler{transactionService: transactionService}
}

// HandleReport handles report endpoints
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetSalesSummary(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetSalesSummary handles GET /api/report endpoints
func (h *ReportHandler) GetSalesSummary(w http.ResponseWriter, r *http.Request) {
	// Check query parameters for date range
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	var summary interface{}
	var err error

	if startDate != "" && endDate != "" {
		// Date range query
		summary, err = h.transactionService.GetSalesSummaryByDateRange(startDate, endDate)
	} else {
		// Today's summary
		summary, err = h.transactionService.GetSalesSummaryToday()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
