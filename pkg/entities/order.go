package entities

type Product struct {
	ID    int
	Name  string
	Price int // in cents
}

type Order struct {
	ID          int
	CustomerID  int
	TotalAmount int
	ProductIDs  []int
}
