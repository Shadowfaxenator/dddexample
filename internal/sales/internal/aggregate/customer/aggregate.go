package customer

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Customer struct {
	aggregate.Aggregate
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func New(name string, age uint, addresses []Address) *Customer {
	return &Customer{
		Aggregate: aggregate.Aggregate{ID: aggregate.NewID()},
		Name:      name,
		Age:       age,
		Addresses: addresses,
	}

}

func (c *Customer) Register() (aggregate.Events[Customer], error) {
	if c.Exists {
		return nil, aggregate.ErrAggregateAlreadyExists
	}
	return aggregate.NewEvents(Registered{
		CustomerID:   c.ID,
		Name:         c.Name,
		Age:          c.Age,
		Addresses:    c.Addresses,
		ActiveOrders: c.ActiveOrders,
	}), nil

}

var ErrInvalidAge = errors.New("invalid age")

func (c *Customer) VerifyOrder(o aggregate.ID) (aggregate.Events[Customer], error) {
	if c.Age < 18 {
		return aggregate.NewEvents(OrderRejected{OrderID: o, Reason: "too young"}), ErrInvalidAge
	}
	return aggregate.NewEvents(OrderAccepted{CustomerID: c.ID, OrderID: o}), nil
}
