package sales

import "ddd/pkg/aggregate"

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	aggregate.ID
	Name         string
	Age          uint
	Addresses    []Address
	ActiveOrders uint
}

// func (p Customer) Events(reg gonvex.RegisterEventFunc[Customer]) {

// 	reg(CustomerCreated{})
// 	reg(CustomerOrderClosed{})
// 	reg(OrderAcceptedByCustomer{})

// }

//agr.Root = NewCustomer("Joe")

// if len(a.cars) >= 10 {
// 	return errors.New("not allowed to add more than 10 cars")
// }

//return agr.Mutate(ctx, gonvex.NewEvent(PERSON_CREATED, BOUNDED_CONTEXT, person))
