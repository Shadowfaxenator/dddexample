package customers

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	ID           domain.CustomerID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

func (c *Customer) Register(cust *Customer) essrv.Events[Customer] {
	if c.ID.IsZero() {
		return essrv.NewEvents(CustomerRegistered{Customer: cust})
	}
	return nil
}

func (c *Customer) VerifyOrder(o domain.OrderID) essrv.Events[Customer] {
	if c.Age < 18 {
		return essrv.NewEvents(OrderRejected{OrderID: o, Reason: "too young"})
	}
	return essrv.NewEvents(OrderAccepted{CustomerID: c.ID, OrderID: o})
}
