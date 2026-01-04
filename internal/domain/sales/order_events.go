package sales

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type OrderCreated struct {
	Order
}

func (ce OrderCreated) Apply(c *Order) {
	*c = ce.Order
}

type CarAddedToOrder struct {
	OrderID eventstore.ID[Order]
	CarID   eventstore.ID[Car]
}

func (ce *CarAddedToOrder) Apply(c *Order) {
	c.Cars[ce.CarID] = struct{}{}
}

type CarRemovedFromOrder struct {
	OrderID eventstore.ID[Order]
	CarID   eventstore.ID[Car]
}

func (ce CarRemovedFromOrder) Apply(c *Order) {
	delete(c.Cars, ce.CarID)
}

type OrderVerified struct {
	OrderID eventstore.ID[Order]
}

func (ce OrderVerified) Apply(c *Order) {
	c.Status = ValidForProcessing
}

type OrderClosed struct{}

func (ce OrderClosed) Apply(c *Order) {
	c.Status = Closed
}
