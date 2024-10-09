package main

import (
	"database/sql"
	"net/http"

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

func initServices(db *sql.DB) *services.OrderService {
	var orderStore services.OrderStore = data.NewOrderDBClient(db)
	return services.NewOrderService(orderStore)
}

func main() {
	db, err := initDB()
	if err != nil {
		return
	}
	r := chi.NewRouter()
	r.Get("/order/{orderID}", func(w http.ResponseWriter, r *http.Request) {
		data.OrderDBClient(db.Conn())
	})
	http.ListenAndServe(":3300", r)
}
