package customer

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type Customer struct {
	ID           aggregate.ID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
	Exists       bool
}

func New(name string, age uint, addresses []Address) *Customer {
	id, err := aggregate.NewID()
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

func (c *Customer) Register(cust *Customer) (aggregate.Events[Customer], error) {
	if c.Exists {
		return nil, eventstore.ErrAggregateAlreadyExists
	}
	return aggregate.NewEvents(&Registered{
		CustomerID:   cust.ID,
		Name:         cust.Name,
		Age:          cust.Age,
		Addresses:    cust.Addresses,
		ActiveOrders: cust.ActiveOrders,
	}), nil

}

var ErrInvalidAge = errors.New("invalid age")

func (c *Customer) VerifyOrder(o aggregate.ID) (aggregate.Events[Customer], error) {

	if c.Age < 18 {
		return aggregate.NewEvents(&OrderRejected{OrderID: o, Reason: "too young"}), ErrInvalidAge
	}
	return aggregate.NewEvents(&OrderAccepted{CustomerID: c.ID, OrderID: o}), nil
}
