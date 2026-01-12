package order

import (
	"slices"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type Posted struct {
	ID         values.OrderID
	CustomerID values.CustomerID
	Cars       []OrderLine
	Status     RentOrderStatus
	Deleted    bool
}

func (ce Posted) Evolve(c *Order) {
	c.ID = ce.ID
	c.Cars = ce.Cars
	c.CustomerID = ce.CustomerID
	c.Status = ce.Status
	c.Deleted = ce.Deleted

}

type CarAdded struct {
	OrderID  values.OrderID
	CarID    values.CarID
	Price    values.Money
	Quantity uint
}

func (ce *CarAdded) Evolve(c *Order) {
	c.Cars = append(c.Cars, OrderLine{CarID: ce.CarID, Price: ce.Price, Quantity: ce.Quantity})
}

type CarRemoved struct {
	OrderID values.OrderID
	CarID   values.CarID
}

func (ce CarRemoved) Evolve(c *Order) {
	c.Cars = slices.DeleteFunc(c.Cars, func(l OrderLine) bool { return l.CarID == ce.CarID })
}

type Verified struct {
	OrderID events.ID[Order]
}

func (ce Verified) Evolve(c *Order) {
	c.Status = StatusValidForProcessing
}

type Closed struct{}

func (ce Closed) Evolve(c *Order) {
	c.Status = StatusClosed
}
