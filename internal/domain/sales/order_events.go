package sales

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type CarAddedToOrder struct {
	OrderID aggregate.ID[Order]
	CarID   aggregate.ID[Car]
}

func (ce *CarAddedToOrder) Apply(c *Order) {
	c.Cars[ce.CarID] = struct{}{}
}

type CarRemovedFromOrder struct {
	OrderID aggregate.ID[Order]
	CarID   aggregate.ID[Car]
}

func (ce *CarRemovedFromOrder) Apply(c *Order) {
	delete(c.Cars, ce.CarID)
}

type OrderVerified struct {
	OrderID aggregate.ID[Order]
}

func (ce *OrderVerified) Apply(c *Order) {
	c.Status = ValidForProcessing
}

type OrderClosed struct {
	OrderID aggregate.ID[Order]
	CustID  aggregate.ID[Customer]
}

func (ce *OrderClosed) Apply(c *Order) {
	c.Status = Closed
}
