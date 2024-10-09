package data

import (
	"database/sql"
	"fmt"

	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type CustomerDBClient struct {
	db *sql.DB
}

func NewCustomerDBClient(db *sql.DB) *CustomerDBClient {
	return &CustomerDBClient{db: db}
}

func (cdb *CustomerDBClient) Get(id int) (e.Customer, error) {
	return e.Customer{}, fmt.Errorf("not implemented")
}

func (cdb *CustomerDBClient) Create(client e.Customer) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
