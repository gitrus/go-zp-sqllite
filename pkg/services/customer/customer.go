package customer

import (
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type Store interface {
	Get(id int) (e.Customer, error)
	Create(customer e.Customer) (e.Customer, error)
}
