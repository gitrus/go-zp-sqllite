package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/data"
	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/services/order"
	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/views"
	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./demo.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db, db.Ping()
}

func InitOrderService(db *sql.DB) *order.OrderService {
	var oStore order.OrderStore = data.NewOrderDBClient(db)
	var pFetcher order.ProductFetcher = data.NewProductDBClient(db)
	return order.NewOrderService(oStore, pFetcher)
}

func main() {
	var db *sql.DB
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatalf("db init failed: %v", err)
		return
	}
	defer db.Close()

	var orderService *order.OrderService = InitOrderService(db)
	var orderViews *views.OrderViews = views.NewOrderViews(*orderService)

	r := chi.NewRouter()
	r.Get("/order/{orderID}", orderViews.GetByID)

	log.Printf("Serve on :3300")
	if err := http.ListenAndServe(":3300", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
