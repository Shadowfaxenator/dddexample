package order

import (
	"slices"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type Posted struct {
	OrderID    aggregate.ID
	CustomerID aggregate.ID
	Cars       OrderLines
	Status     RentOrderStatus
	Total      values.Money
}

func (ce *Posted) Evolve(c *Order) {
	c.ID = ce.OrderID
	c.Cars = ce.Cars
	c.CustomerID = ce.CustomerID
	c.Status = ce.Status
	c.Total = ce.Total

}

type CarAdded struct {
	OrderID  aggregate.ID
	CarID    aggregate.ID
	Price    values.Money
	Quantity uint
}

func (ce *CarAdded) Evolve(c *Order) {
	c.Cars = append(c.Cars, OrderLine{CarID: ce.CarID, Price: ce.Price, Quantity: ce.Quantity})
}

type CarRemoved struct {
	OrderID aggregate.ID
	CarID   aggregate.ID
}

func (ce *CarRemoved) Evolve(c *Order) {
	c.Cars = slices.DeleteFunc(c.Cars, func(l OrderLine) bool { return l.CarID == ce.CarID })
}

type Verified struct {
	OrderID aggregate.ID
}

func (ce *Verified) Evolve(c *Order) {
	c.Status = StatusValidForProcessing
}

type Closed struct{}

func (ce *Closed) Evolve(c *Order) {
	c.Status = StatusClosed
}
