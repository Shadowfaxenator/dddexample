package customer

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	eventstore1 "github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Customer struct {
	ID           eventstore1.ID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
	Exists       bool
}

func New(name string, age uint, addresses []Address) *Customer {
	id, err := eventstore1.NewID()
	if err != nil {
		panic(err)
	}

	return &Customer{
		ID:        id,
		Name:      name,
		Age:       age,
		Addresses: addresses,
	}

}

func (c *Customer) Register(cust *Customer) (eventstore1.Events[Customer], error) {
	if c.Exists {
		return nil, aggregate.ErrAggregateAlreadyExists
	}
	return eventstore1.NewEvents(&Registered{
		CustomerID:   cust.ID,
		Name:         cust.Name,
		Age:          cust.Age,
		Addresses:    cust.Addresses,
		ActiveOrders: cust.ActiveOrders,
	}), nil

}

var ErrInvalidAge = errors.New("invalid age")

func (c *Customer) VerifyOrder(o eventstore1.ID) (eventstore1.Events[Customer], error) {

	if c.Age < 18 {
		return eventstore1.NewEvents(&OrderRejected{OrderID: o, Reason: "too young"}), ErrInvalidAge
	}
	return eventstore1.NewEvents(&OrderAccepted{CustomerID: c.ID, OrderID: o}), nil
}
