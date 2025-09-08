package sales

import (
	"context"
	"ddd/pkg/aggregate"
)

type CreateOrder struct {
	OrderID aggregate.ID
	CustID  aggregate.ID
}

func (c CreateOrder) Execute(o *Order) (*aggregate.Event[Order], error) {
	event := aggregate.NewEvent(OrderCreated{
		Order{ID: c.OrderID, CustomerID: c.CustID,
			Cars: make(map[aggregate.ID]struct{}), Status: ProcessingByCustomer,
		}})
	return event, nil
}

func (c *SubDomain) CreateOrder(ctx context.Context, ordid aggregate.ID, custID aggregate.ID) error {

	return c.order.Command(ctx, ordid, CreateOrder{CustID: custID, OrderID: ordid})
}

func (c *SubDomain) CloseOrder(ctx context.Context, orderID aggregate.ID) error {
	return c.order.CommandFunc(ctx, orderID, func(c *Order) (*aggregate.Event[Order], error) {
		return aggregate.NewEvent(OrderClosed{OrderID: orderID, CustID: c.CustomerID}), nil
	})
}

func (c *SubDomain) AddCarToOrder(ctx context.Context, orderID aggregate.ID, carID aggregate.ID) error {
	return c.order.CommandFunc(ctx, orderID, func(c *Order) (*aggregate.Event[Order], error) {
		//return &CustomerCreated{}, nil
		return aggregate.NewEvent(CarAddedToOrder{CarID: carID, OrderID: orderID}), nil
	})
}

func (c *SubDomain) MarkVerrified(ctx context.Context, orderID aggregate.ID) error {
	return c.order.CommandFunc(ctx, orderID, func(c *Order) (*aggregate.Event[Order], error) {
		return aggregate.NewEvent(OrderVerified{}), nil
	})
}
