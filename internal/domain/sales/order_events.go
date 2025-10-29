package sales

import (
	"ddd/pkg/aggregate"
)

type OrderCreated struct {
	Order Order
}

func (ce OrderCreated) Apply(c *Order) {
	*c = ce.Order
}

func (ce OrderCreated) String() string {
	return "ORDER_CREATED"
}

type CarAddedToOrder struct {
	OrderID aggregate.ID[Order]
	CarID   aggregate.ID[Car]
}

func (ce CarAddedToOrder) Apply(c *Order) {
	c.Cars[ce.CarID] = struct{}{}
}

type CarRemovedFromOrder struct {
	CarID aggregate.ID[Car]
}

func (ce CarRemovedFromOrder) Apply(c *Order) {
	delete(c.Cars, ce.CarID)
}

type OrderVerified struct {
}

func (ce OrderVerified) Apply(c *Order) {
	c.Status = ValidForProcessing
}

type OrderClosed struct {
	OrderID aggregate.ID[Order]
	CustID  aggregate.ID[Customer]
}

func (ce OrderClosed) Apply(c *Order) {
	c.Status = Closed
}
