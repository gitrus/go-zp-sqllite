package views

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/services/order"
	"github.com/go-chi/chi/v5"
)

type OrderViews struct {
	orderService order.OrderService
}

func NewOrderViews(os order.OrderService) *OrderViews {
	return &OrderViews{orderService: os}
}

func (ov *OrderViews) GetByID(w http.ResponseWriter, r *http.Request) {
	orderIDStr := chi.URLParam(r, "orderID")
	if orderIDStr == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := ov.orderService.GetByID(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode order", http.StatusInternalServerError)
	}
}
