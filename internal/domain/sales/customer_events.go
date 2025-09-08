package sales

import (
	"ddd/pkg/aggregate"
)

type CustomerCreated struct {
	Customer Customer
}

func (cc CustomerCreated) Apply(c *Customer) {
	*c = cc.Customer
}

type CustomerOrderClosed struct {
	OrderID aggregate.ID
}

func (CustomerOrderClosed) Apply(c *Customer) {
	c.ActiveOrders--
}

type OrderAccepted struct {
	OrderID aggregate.ID
}

func (OrderAccepted) Apply(c *Customer) {
	c.ActiveOrders++
}
