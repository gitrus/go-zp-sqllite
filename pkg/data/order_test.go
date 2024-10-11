package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/data"
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type OrderDBClientTestSuite struct {
	suite.Suite
	db     *sql.DB
	client *data.OrderDBClient
}

func (suite *OrderDBClientTestSuite) SetupSuite() {
	var err error
	// arrange
	suite.db, err = sql.Open("sqlite", ":memory:")
	assert.NoError(suite.T(), err, "Failed to open test database")

	// arrange
	suite.createTables()
}

func (suite *OrderDBClientTestSuite) SetupTest() {
	// teardown
	suite.eraseDB()

	suite.client = data.NewOrderDBClient(suite.db)
}

func (suite *OrderDBClientTestSuite) TearDownTest() {
	suite.eraseDB()
}

func (suite *OrderDBClientTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *OrderDBClientTestSuite) eraseDB() {
	_, err := suite.db.Exec("DELETE FROM order_products")
	assert.NoError(suite.T(), err, "Failed to clear order_products table in teardown")

	_, err = suite.db.Exec("DELETE FROM orders")
	assert.NoError(suite.T(), err, "Failed to clear orders table in teardown")
}

func (suite *OrderDBClientTestSuite) createTables() {
	orderTable := `
    CREATE TABLE orders (
        order_id INTEGER PRIMARY KEY AUTOINCREMENT,
        customer_id INTEGER NOT NULL CHECK (customer_id > 0),
        total_amount INTEGER NOT NULL
    );`

	orderProductsTable := `
    CREATE TABLE order_products (
        order_id INTEGER NOT NULL,
        product_id INTEGER NOT NULL,
        FOREIGN KEY(order_id) REFERENCES orders(order_id)
    );`

	_, err := suite.db.Exec(orderTable)
	assert.NoError(suite.T(), err, "Failed to create orders table")

	_, err = suite.db.Exec(orderProductsTable)
	assert.NoError(suite.T(), err, "Failed to create order_products table")
}

func (suite *OrderDBClientTestSuite) TestOrderDBClient_Create_and_Get() {
	testCases := []struct {
		name             string
		customerID       int
		productIDs       []int
		orderTotalAmount int
		expectedError    error
	}{
		{
			name:             "Successful creation with products",
			customerID:       101,
			productIDs:       []int{1, 2, 3},
			orderTotalAmount: 600,
			expectedError:    nil,
		},
		{
			name:             "Successful creation without products",
			customerID:       102,
			productIDs:       []int{},
			orderTotalAmount: 400,
			expectedError:    nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// act
			orderID, err := suite.client.Create(tc.customerID, tc.productIDs, tc.orderTotalAmount)

			if tc.expectedError != nil {
				// assert
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), 0, orderID)
			} else {
				// assert
				assert.NoError(suite.T(), err)
				assert.NotEqual(suite.T(), 0, orderID)

				// act (act after assert is not best practice)
				order, err := suite.client.Get(orderID)
				assert.NoError(suite.T(), err)

				// assert
				assert.Equal(suite.T(), tc.customerID, order.CustomerID)
				assert.Equal(suite.T(), len(tc.productIDs), len(order.ProductIDs))
				if len(tc.productIDs) > 0 {
					assert.Equal(suite.T(), tc.productIDs, order.ProductIDs)
				}
			}
		})
	}
}

func (suite *OrderDBClientTestSuite) TestOrderDBClient_Get() {
	// arrange
	orderID1, err := suite.client.Create(100, []int{10, 20, 30}, 500)
	assert.NoError(suite.T(), err, "Failed to insert order 1")

	testCases := []struct {
		name          string
		orderID       int
		expectedOrder e.Order
		expectedError error
	}{
		{
			name:    "Order exists with products",
			orderID: orderID1,
			expectedOrder: e.Order{
				ID:         orderID1,
				CustomerID: 100,
				ProductIDs: []int{10, 20, 30},
			},
			expectedError: nil,
		},
		{
			name:          "Order does not exist",
			orderID:       999,
			expectedOrder: e.Order{},
			expectedError: fmt.Errorf("order with ID %d not found", 999),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// act
			order, err := suite.client.Get(tc.orderID)

			if tc.expectedError != nil {
				// assert
				assert.Error(suite.T(), err)
				assert.EqualError(suite.T(), err, tc.expectedError.Error())
				assert.Empty(suite.T(), order.ProductIDs)
				assert.Equal(suite.T(), e.Order{ID: 0, CustomerID: 0, TotalAmount: 0, ProductIDs: []int(nil)}, order)
			} else {
				// assert
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedOrder, order)
			}
		})
	}
}

func TestOrderDBClientTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDBClientTestSuite))
}
