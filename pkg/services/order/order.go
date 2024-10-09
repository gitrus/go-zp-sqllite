package order

import (
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/services/customer"
)

type Store interface {
	Get(id int) (e.Order, error)
	Create(customerID int, productIDs []int) (int, error)
}

type OrderService struct {
	orderStore    Store
	customerStore customer.Store
}

func NewOrderService(customerStore customer.Store, orderStore Store) *OrderService {
	return &OrderService{customerStore: customerStore, orderStore: orderStore}
}
