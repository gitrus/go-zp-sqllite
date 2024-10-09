package data

import (
	"database/sql"
	"fmt"

	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type OrderDBClient struct {
	db *sql.DB
}

func NewOrderDBClient(db *sql.DB) *OrderDBClient {
	return &OrderDBClient{db: db}
}

func (cdb *OrderDBClient) Get(id int) (e.Order, error) {
	var order e.Order

	row := cdb.db.QueryRow("SELECT id, customer_id FROM orders WHERE id = :id", sql.Named("id", id))
	err := row.Scan(&order.ID, &order.CustomerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("order with ID %d not found", id)
		}
		return order, err
	}

	rows, err := cdb.db.Query("SELECT product_id FROM order_products WHERE order_id = :id", sql.Named("id", id))
	if err != nil {
		return order, err
	}
	defer rows.Close()

	var productID int
	for rows.Next() {
		if err := rows.Scan(&productID); err != nil {
			return order, err
		}
		order.ProductIDs = append(order.ProductIDs, productID)
	}

	if err = rows.Err(); err != nil {
		return order, err
	}

	return order, nil
}

func (cdb *OrderDBClient) Create(customerID int, productIDs []int) (int, error) {
	// Start a transaction
	tx, err := cdb.db.Begin()
	if err != nil {
		return 0, err
	}

	// Rollback the transaction on error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.Exec(
		"INSERT INTO orders (customer_id) VALUES (:customerID)",
		sql.Named("customerID", customerID),
	)
	if err != nil {
		return 0, err
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, productID := range productIDs {
		_, err = tx.Exec(
			"INSERT INTO order_products (order_id, product_id) VALUES (:orderID, :productID)",
			sql.Named("orderID", orderID),
			sql.Named("productID", productID),
		)
		if err != nil {
			return 0, err
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return int(orderID), nil
}
