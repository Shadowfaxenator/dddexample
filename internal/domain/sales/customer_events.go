package sales

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type CustomerEvents = eventstore.Events[Customer]

type CustomerCreated struct {
	Customer
}

func (e CustomerCreated) Apply(c *Customer) {
	*c = e.Customer
}

type CustomerOrderClosed struct {
	CustomerID eventstore.ID[Customer]
	OrderID    eventstore.ID[Order]
}

func (CustomerOrderClosed) Apply(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	OrderID eventstore.ID[Order]
}

func (OrderAccepted) Apply(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	CustomerID eventstore.ID[Customer]
	OrderID    eventstore.ID[Order]
	Error      string
}

func (OrderRejected) Apply(c *Customer) {

}
