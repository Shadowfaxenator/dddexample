package sales

import (
	"context"
	"ddd/pkg/aggregate"
	"fmt"

	"errors"
)

type CreateCustomer struct {
	Customer
}

func (c CreateCustomer) Execute(a *Customer) (*aggregate.Event[Customer], error) {

	//aggregate.WithType(CustomerCreated{Customer: c.Customer})
	return aggregate.NewEvent(CustomerCreated{Customer: c.Customer}), nil
}

func (c *SubDomain) removeFromOneOrder(ctx context.Context, custid aggregate.ID) error {
	return c.customer.CommandFunc(ctx, custid, func(c *Customer) (*aggregate.Event[Customer], error) {
		if c.ActiveOrders <= 0 {
			return nil, ErrMinOrders
		}
		return aggregate.NewEvent(CustomerOrderClosed{OrderID: aggregate.NewID()}), nil
	})
}

type ValidateOrdersError struct {
	etype string
	value int
}

func (e ValidateOrdersError) Error() string {
	return fmt.Sprintf("active orders %s %d", e.etype, e.value)
}

var ErrMaxOrders = &ValidateOrdersError{">=", 3}
var ErrMinOrders = &ValidateOrdersError{"<", 0}
var ErrValidateAge = errors.New("age is < 18")

func (c *SubDomain) ValidateOrder(ctx context.Context, custid aggregate.ID, orderID aggregate.ID) error {
	return c.customer.CommandFunc(ctx, custid, func(c *Customer) (*aggregate.Event[Customer], error) {

		switch {
		case c.Age <= 18:
			return nil, ErrValidateAge
		case c.ActiveOrders >= 3:
			return nil, ErrMaxOrders
		}

		return aggregate.NewEvent(OrderAccepted{OrderID: orderID}), nil
	})
}

func (c *SubDomain) CreateCustomer(ctx context.Context, id aggregate.ID, name string, age uint) error {

	return c.customer.Command(ctx, id, CreateCustomer{Customer: Customer{ID: id, Name: name, Age: age}})
	// return c.customer.CommandFunc(ctx, nil, func(c *Customer) (aggregate.Event[Customer], error) {
	// 	return &CustomerCreated{Customer: Customer{Name: name, Age: age, ActiveOrders: 0}}, nil
	// })
}
