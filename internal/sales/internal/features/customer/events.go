package customer

import (
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
)

type Registered struct {
	ID           ids.CustomerID
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
	CustomerID ids.CustomerID
	OrderID    ids.OrderID
}

func (OrderClosed) Evolve(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	CustomerID ids.CustomerID
	OrderID    ids.OrderID
}

func (OrderAccepted) Evolve(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID ids.OrderID
	Reason  string
}

func (OrderRejected) Evolve(c *Customer) {

}
