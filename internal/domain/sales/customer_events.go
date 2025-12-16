package sales

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type CustomerOrderClosed struct {
	CustomerID aggregate.ID[Customer]
	OrderID    aggregate.ID[Order]
}

func (*CustomerOrderClosed) Apply(c *Customer) {
	c.ActiveOrders--

}

type OrderAccepted struct {
	OrderID aggregate.ID[Order]
}

func (*OrderAccepted) Apply(c *Customer) {
	c.ActiveOrders++
}

type OrderRejected struct {
	CustomerID aggregate.ID[Customer]
	OrderID    aggregate.ID[Order]
	Error      string
}

func (*OrderRejected) Apply(c *Customer) {

}
