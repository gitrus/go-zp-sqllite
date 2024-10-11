// order_db_client_test.go
package data_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/data"
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite" // Import SQLite driver for side-effects
)

// OrderDBClientTestSuite defines the test suite for OrderDBClient
type OrderDBClientTestSuite struct {
	suite.Suite
	db     *sql.DB
	tx     *sql.Tx
	client *data.OrderDBClient
}

// SetupSuite runs once before the suite starts
func (suite *OrderDBClientTestSuite) SetupSuite() {
	// Initialize in-memory SQLite database
	var err error
	suite.db, err = sql.Open("sqlite", ":memory:")
	assert.NoError(suite.T(), err, "Failed to open test database")

	// Create necessary tables with constraints
	suite.createTables()
}

// SetupTest runs before each test in the suite
func (suite *OrderDBClientTestSuite) SetupTest() {
	// Begin a new transaction
	var err error
	suite.tx, err = suite.db.Begin()
	assert.NoError(suite.T(), err, "Failed to begin transaction")

	// Assign the transaction to OrderDBClient
	suite.client = data.NewOrderDBClient(suite.tx)
}

// TearDownTest runs after each test in the suite
func (suite *OrderDBClientTestSuite) TearDownTest() {
	// Rollback the transaction to revert changes
	err := suite.tx.Rollback()
	assert.NoError(suite.T(), err, "Failed to rollback transaction")
}

// TearDownSuite runs once after the suite finishes
func (suite *OrderDBClientTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// createTables creates the necessary tables for testing with constraints
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

// TestOrderDBClient_Create tests the Create method of OrderDBClient
func (suite *OrderDBClientTestSuite) TestOrderDBClient_Create() {
	testCases := []struct {
		name             string
		customerID       int
		productIDs       []int
		orderTotalAmount int
		expectedID       int
		expectedError    error
	}{
		{
			name:             "Successful creation with products",
			customerID:       101,
			productIDs:       []int{1, 2, 3},
			orderTotalAmount: 600,
			expectedID:       1, // First auto-incremented ID within transaction
			expectedError:    nil,
		},
		{
			name:             "Successful creation without products",
			customerID:       102,
			productIDs:       []int{},
			orderTotalAmount: 400,
			expectedID:       1, // Reset to 1 for isolated test
			expectedError:    nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			orderID, err := suite.client.Create(tc.customerID, tc.productIDs, tc.orderTotalAmount)

			if tc.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), 0, orderID)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedID, orderID)

				// Verify that the order exists in the database
				order, err := suite.client.Get(orderID)
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.customerID, order.CustomerID)
				assert.Equal(suite.T(), tc.productIDs, order.ProductIDs)
			}
		})
	}
}

// TestOrderDBClient_Get tests the Get method of OrderDBClient
func (suite *OrderDBClientTestSuite) TestOrderDBClient_Get() {
	// Pre-insert some orders within the transaction
	orderID1, err := suite.client.Create(100, []int{10, 20, 30}, 500)
	assert.NoError(suite.T(), err, "Failed to insert order 1")

	orderID2, err := suite.client.Create(200, []int{}, 300)
	assert.NoError(suite.T(), err, "Failed to insert order 2")

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
			name:    "Order exists without products",
			orderID: orderID2,
			expectedOrder: e.Order{
				ID:         orderID2,
				CustomerID: 200,
				ProductIDs: []int{},
			},
			expectedError: nil,
		},
		{
			name:          "Order does not exist",
			orderID:       999, // Assuming this ID does not exist
			expectedOrder: e.Order{},
			expectedError: fmt.Errorf("order with ID %d not found", 999),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			order, err := suite.client.Get(tc.orderID)

			if tc.expectedError != nil {
				assert.Error(suite.T(), err)
				assert.EqualError(suite.T(), err, tc.expectedError.Error())
				assert.Equal(suite.T(), e.Order{}, order)
			} else {
				assert.NoError(suite.T(), err)
				assert.Equal(suite.T(), tc.expectedOrder, order)
			}
		})
	}
}

// TestOrderDBClient_NotFound ensures that Get returns an error for non-existent orders
func (suite *OrderDBClientTestSuite) TestOrderDBClient_NotFound() {
	orderID := 999 // Non-existent ID
	order, err := suite.client.Get(orderID)
	assert.Error(suite.T(), err)
	assert.EqualError(suite.T(), err, fmt.Sprintf("order with ID %d not found", orderID))
	assert.Equal(suite.T(), e.Order{}, order)
}

// TestOrderDBClient_InvalidCreate tests Create with invalid data
func (suite *OrderDBClientTestSuite) TestOrderDBClient_InvalidCreate() {
	// Example: Negative customerID, assuming business logic disallows it
	customerID := -1
	productIDs := []int{1, 2}
	orderTotalAmount := 100

	orderID, err := suite.client.Create(customerID, productIDs, orderTotalAmount)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), 0, orderID)
}

// In order to run the suite, we need a Test function
func TestOrderDBClientTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDBClientTestSuite))
}
