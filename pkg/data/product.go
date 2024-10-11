package data

import (
	"database/sql"
	"fmt"
	"strings"

	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

// ProductDBClient provides an implementation for ProductFetcher
type ProductDBClient struct {
	db *sql.DB
}

func NewProductDBClient(db *sql.DB) *ProductDBClient {
	return &ProductDBClient{db: db}
}

func (pdb *ProductDBClient) GetMultiple(ids []int) ([]e.Product, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma
	query := fmt.Sprintf("SELECT id, name, price FROM products WHERE id IN (%s)", placeholders)

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = sql.Named(fmt.Sprintf("id%d", i+1), id)
	}

	rows, err := pdb.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []e.Product
	for rows.Next() {
		var product e.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pdb *ProductDBClient) Create(product e.Product) (int, error) {
	result, err := pdb.db.Exec("INSERT INTO products (name, price) VALUES (:name, :price)",
		sql.Named("name", product.Name),
		sql.Named("price", product.Price))
	if err != nil {
		return 0, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(productID), nil
}
