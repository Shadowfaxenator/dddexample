package sales

import (
	"ddd/pkg/aggregate"
	"fmt"
)

type CreateCustomer struct {
	Customer
}

func (c CreateCustomer) Execute(a *Customer) aggregate.Event[Customer] {
	if a != nil {
		return aggregate.EventError[Customer]{Reason: "customer already exists"}
	}
	//ddd.WithType(CustomerCreated{Customer: c.Customer})
	return &CustomerCreated{Customer: c.Customer}
}

type ValidateOrdersError struct {
	etype string
	value int
}

func (e ValidateOrdersError) Error() string {
	return fmt.Sprintf("active orders %s %d", e.etype, e.value)
}

type ValidateAgeError struct {
	age uint
}

func (e ValidateAgeError) Error() string {
	return fmt.Sprintf("age is %d, must be greater than 18", e.age)
}

var ErrMaxOrders = &ValidateOrdersError{">=", 3}
var ErrMinOrders = &ValidateOrdersError{"<", 0}

func NewValidateAgeError(age uint) *ValidateAgeError {
	return &ValidateAgeError{age}
}

type ValidateOrder struct {
	OrderID aggregate.ID[Order]
}

func (v ValidateOrder) Execute(c *Customer) aggregate.Event[Customer] {
	if c.Age <= 18 {
		return &OrderRejected{OrderID: v.OrderID, Error: NewValidateAgeError(c.Age)}
	}
	if c.ActiveOrders >= 3 {
		return &OrderRejected{OrderID: v.OrderID, Error: ErrMaxOrders}
	}

	return &OrderAccepted{OrderID: v.OrderID}
}
