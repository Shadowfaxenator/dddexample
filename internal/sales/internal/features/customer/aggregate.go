package customer

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
)

type Customer struct {
	ID           ids.CustomerID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func (c *Customer) Register() (events.Events[Customer], error) {

	return events.New(Registered{
		ID:           c.ID,
		Name:         c.Name,
		Age:          c.Age,
		Addresses:    c.Addresses,
		ActiveOrders: c.ActiveOrders,
	}), nil

}

var ErrInvalidAge = errors.New("invalid age")

func (c *Customer) VerifyOrder(o ids.OrderID) (events.Events[Customer], error) {
	if c.Age < 18 {
		return events.New(OrderRejected{OrderID: o, Reason: "too young"}), ErrInvalidAge
	}
	return events.New(OrderAccepted{CustomerID: c.ID, OrderID: o}), nil
}
