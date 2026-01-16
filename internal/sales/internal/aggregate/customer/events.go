package customer

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Registered struct {
	CustomerID   aggregate.ID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func (e *Registered) Evolve(c *Customer) {
	c.Exists = true
	c.ID = e.CustomerID
	c.Name = e.Name
	c.Age = e.Age
	c.Addresses = e.Addresses
	c.ActiveOrders = e.ActiveOrders
}

type OrderClosed struct {
	CustomerID aggregate.ID
	OrderID    aggregate.ID
}

func (e *OrderClosed) Evolve(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	CustomerID aggregate.ID
	OrderID    aggregate.ID
}

func (e *OrderAccepted) Evolve(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	OrderID aggregate.ID
	Reason  string
}

func (e *OrderRejected) Evolve(c *Customer) {

}
