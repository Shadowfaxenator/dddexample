package order

import (
	"slices"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/money"
)

type Posted struct {
	ID         ids.OrderID
	CustomerID ids.CustomerID
	Cars       []OrderLine
	Status     RentOrderStatus
	Deleted    bool
}

func (ce Posted) Evolve(c *AggregateRoot) {
	c.ID = ce.ID
	c.Cars = ce.Cars
	c.CustomerID = ce.CustomerID
	c.Status = ce.Status
	c.Deleted = ce.Deleted

}

type CarAdded struct {
	OrderID  ids.OrderID
	CarID    ids.CarID
	Price    money.Money
	Quantity uint
}

func (ce *CarAdded) Evolve(c *AggregateRoot) {
	c.Cars = append(c.Cars, OrderLine{CarID: ce.CarID, Price: ce.Price, Quantity: ce.Quantity})
}

type CarRemoved struct {
	OrderID ids.OrderID
	CarID   ids.CarID
}

func (ce CarRemoved) Evolve(c *AggregateRoot) {
	c.Cars = slices.DeleteFunc(c.Cars, func(l OrderLine) bool { return l.CarID == ce.CarID })
}

type Verified struct {
	OrderID essrv.ID[AggregateRoot]
}

func (ce Verified) Evolve(c *AggregateRoot) {
	c.Status = StatusValidForProcessing
}

type Closed struct{}

func (ce Closed) Evolve(c *AggregateRoot) {
	c.Status = StatusClosed
}
