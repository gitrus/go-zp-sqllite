package services

import (
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
)

type OrderStore interface {
	Get(id int) (e.Order, error)
	Create(customerID int, productIDs []int, orderTotalAmount int) (int, error)
}

type ProductFetcher interface {
	GetMultiple(ids []int) ([]e.Product, error)
}

type OrderService struct {
	orderStore     OrderStore
	productFetcher ProductFetcher
}

func NewOrderService(orderStore OrderStore) *OrderService {
	return &OrderService{orderStore: orderStore}
}

func (os *OrderService) GetByID(id int) (e.Order, error) {
	return os.orderStore.Get(id)
}

func (os *OrderService) Create(customerID int, productIDs []int) (e.Order, error) {
	var err error
	var products []e.Product
	products, err = os.productFetcher.GetMultiple(productIDs)
	if err != nil {
		return e.Order{}, nil
	}

	var orderTotalAmount int
	for _, product := range products {
		orderTotalAmount += product.Price
	}

	var orderID int

	orderID, err = os.orderStore.Create(customerID, productIDs, orderTotalAmount)
	if err != nil {
		return e.Order{}, err
	}
	return os.GetByID(orderID)
}
