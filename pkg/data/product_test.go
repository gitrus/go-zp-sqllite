package data_test

import (
	"database/sql"
	"testing"

	"github.com/Yandex-Practicum/go-db-sql-query-test/pkg/data"
	e "github.com/Yandex-Practicum/go-db-sql-query-test/pkg/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type ProductDBClientTestSuite struct {
	suite.Suite
	db     *sql.DB
	client *data.ProductDBClient
}

func (suite *ProductDBClientTestSuite) SetupSuite() {
	var err error
	suite.db, err = sql.Open("sqlite", ":memory:")
	assert.NoError(suite.T(), err, "Failed to open test database")

	suite.createTables()
	suite.client = data.NewProductDBClient(suite.db)
}

func (suite *ProductDBClientTestSuite) createTables() {
	productTable := `
    CREATE TABLE products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        price INTEGER NOT NULL
    );`

	_, err := suite.db.Exec(productTable)
	assert.NoError(suite.T(), err, "Failed to create products table")
}

func (suite *ProductDBClientTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *ProductDBClientTestSuite) SetupTest() {
	suite.eraseDB()
}

func (suite *ProductDBClientTestSuite) eraseDB() {
	_, err := suite.db.Exec("DELETE FROM products")
	assert.NoError(suite.T(), err, "Failed to clear products table")
}

func (suite *ProductDBClientTestSuite) TestProductDBClient_Create_and_GetMultiple() {
	// arrange
	product1 := e.Product{Name: "Product A", Price: 1000}
	product2 := e.Product{Name: "Product B", Price: 2000}

	productID1, err := suite.client.Create(product1)
	assert.NoError(suite.T(), err, "Failed to create product 1")
	assert.NotEqual(suite.T(), 0, productID1)

	productID2, err := suite.client.Create(product2)
	assert.NoError(suite.T(), err, "Failed to create product 2")
	assert.NotEqual(suite.T(), 0, productID2)

	// act
	productIDs := []int{productID1, productID2}
	products, err := suite.client.GetMultiple(productIDs)

	// assert
	assert.NoError(suite.T(), err, "Failed to get multiple products")
	assert.Len(suite.T(), products, 2)

	assert.Equal(suite.T(), productID1, products[0].ID)
	assert.Equal(suite.T(), "Product A", products[0].Name)
	assert.Equal(suite.T(), 1000, products[0].Price)

	assert.Equal(suite.T(), productID2, products[1].ID)
	assert.Equal(suite.T(), "Product B", products[1].Name)
	assert.Equal(suite.T(), 2000, products[1].Price)
}

func TestProductDBClientTestSuite(t *testing.T) {
	suite.Run(t, new(ProductDBClientTestSuite))
}
