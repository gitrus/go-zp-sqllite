package services

import (
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type OrderStore interface {
	Get(id int) (e.Order, error)
	Create(userId int, products []e.Product) (e.Order, error)
}

type OrderService struct {
	orderStore OrderStore
}

func NewOrderService(orderStore OrderStore) *OrderService {
	return &OrderService{orderStore: orderStore}
}
