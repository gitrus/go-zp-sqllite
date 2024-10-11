package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/data"
	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/services"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./demo.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, db.Ping()
}

// logic
func initServices(db *sql.DB) *services.OrderService {
	var orderStore services.OrderStore = data.NewOrderDBClient(db)
	return services.NewOrderService(orderStore)
}

func main() {
	var db *sql.DB
	var err error
	db, err = initDB()
	if err != nil {
		return
	}
	defer db.Close()

	var orderService *services.OrderService = initServices(db)

	r := chi.NewRouter()

	r.Get("/order/{orderID}", func(w http.ResponseWriter, r *http.Request) {
		orderIDStr := chi.URLParam(r, "orderID")
		orderID, err := strconv.Atoi(orderIDStr)
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		order, err := orderService.GetByID(orderID)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "Failed to encode order", http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":3300", r)
}
