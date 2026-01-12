package customer

import "github.com/alekseev-bro/dddexample/internal/sales/internal/values"

type Registered struct {
	ID           values.CustomerID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func (e Registered) Evolve(c *Customer) {
	c.ID = e.ID
	c.Name = e.Name
	c.Age = e.Age
	c.Addresses = e.Addresses
	c.ActiveOrders = e.ActiveOrders
}

type OrderClosed struct {
	CustomerID values.CustomerID
	OrderID    values.OrderID
}

func (OrderClosed) Evolve(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	CustomerID values.CustomerID
	OrderID    values.OrderID
}

func (OrderAccepted) Evolve(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID values.OrderID
	Reason  string
}

func (OrderRejected) Evolve(c *Customer) {

}
