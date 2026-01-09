package orders

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain"
)

type OrderPosted struct {
	*Order
}

func (ce OrderPosted) Evolve(c *Order) {
	*c = *ce.Order
}

type CarAddedToOrder struct {
	OrderID domain.OrderID
	CarID   domain.CarID
}

func (ce *CarAddedToOrder) Evolve(c *Order) {
	c.Cars[ce.CarID] = struct{}{}
}

type CarRemovedFromOrder struct {
	OrderID domain.OrderID
	CarID   domain.CarID
}

func (ce CarRemovedFromOrder) Evolve(c *Order) {
	delete(c.Cars, ce.CarID)
}

type OrderVerified struct {
	OrderID essrv.ID[Order]
}

func (ce OrderVerified) Evolve(c *Order) {
	c.Status = ValidForProcessing
}

type OrderClosed struct{}

func (ce OrderClosed) Evolve(c *Order) {
	c.Status = Closed
}
