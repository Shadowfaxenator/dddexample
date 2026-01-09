package customer

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
)

type AggregateRoot struct {
	ID           ids.CustomerID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func (c *AggregateRoot) Register() (essrv.Events[AggregateRoot], error) {

	return essrv.NewEvents(Registered{
		ID:           c.ID,
		Name:         c.Name,
		Age:          c.Age,
		Addresses:    c.Addresses,
		ActiveOrders: c.ActiveOrders,
	}), nil

}

var ErrInvalidAge = errors.New("invalid age")

func (c *AggregateRoot) VerifyOrder(o ids.OrderID) (essrv.Events[AggregateRoot], error) {
	if c.Age < 18 {
		return essrv.NewEvents(OrderRejected{OrderID: o, Reason: "too young"}), ErrInvalidAge
	}
	return essrv.NewEvents(OrderAccepted{CustomerID: c.ID, OrderID: o}), nil
}
