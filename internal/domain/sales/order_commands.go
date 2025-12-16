package sales

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

// type CreateOrder struct {
// 	OrderID aggregate.ID[Order]
// 	CustID  aggregate.ID[Customer]
// }

// func (c *CreateOrder) Execute(o *Order) aggregate.Event[Order] {
// 	event := &OrderCreated{
// 		Order{ID: c.OrderID, CustomerID: c.CustID,
// 			Cars: make(map[aggregate.ID[Car]]struct{}), Status: ProcessingByCustomer,
// 		}}
// 	return event
// }

// func (e *CreateOrder) AggregateID() aggregate.ID[Order] {
// 	return e.OrderID
// }

type CloseOrder struct {
	OrderID aggregate.ID[Order]
	CustID  aggregate.ID[Customer]
}

func (c *CloseOrder) Execute(o *Order) aggregate.Event[Order] {
	event := &OrderClosed{
		OrderID: c.OrderID,
		CustID:  o.CustomerID,
	}
	return event
}

func (c *CloseOrder) AggregateID() aggregate.ID[Order] {
	return c.OrderID
}

type AddCarToOrder struct {
	OrderID aggregate.ID[Order]
	CarID   aggregate.ID[Car]
}

func (c *AddCarToOrder) Execute(o *Order) aggregate.Event[Order] {
	event := &CarAddedToOrder{
		OrderID: c.OrderID,
		CarID:   c.CarID,
	}
	return event
}

func (c *AddCarToOrder) AggregateID() aggregate.ID[Order] {
	return c.OrderID
}
