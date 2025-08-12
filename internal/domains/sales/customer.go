package sales

import (
	"context"
	"errors"
	"ttt/gonvex"
)

type CustomerAggregate struct {
	Stream gonvex.EventStream[*Customer]
	Create gonvex.Command[CustomerCreated, *Customer]
	AddCar gonvex.Command[CustomerCarAdded, *Customer]
}

func NewCustomer(name string) *Customer {
	return &Customer{Name: name, Cars: make([]gonvex.AggregateID, 0)}
}
func NewCustomerAggregate(ctx context.Context) *CustomerAggregate {
	return &CustomerAggregate{
		Stream: gonvex.NewEventStream[Customer](ctx),
		Create: gonvex.NewCommand[CustomerCreated](func(a *Customer) error { return nil }),
		AddCar: gonvex.NewCommand[CustomerCarAdded](func(a *Customer) error {
			if len(a.Cars) >= 3 {
				return errors.New("not allowed to add more than 10 cars")
			}
			return nil
		}),
	}
}

type Customer struct {
	Name string
	Cars []gonvex.AggregateID
}

func (p *Customer) Reduce(event any) error {
	//	fmt.Printf("ev: %+v\n", event)
	switch ev := event.(type) {
	case CustomerCarAdded:

		id := gonvex.AggregateID(ev)
		p.Cars = append(p.Cars, id)
	case CustomerCreated:

		per := Customer(ev)

		*p = per

	}

	return nil
}

// func (c *CustomerAggregate) Create(ctx context.Context) (gonvex.AggregateID, error) {
// 	agr := c.Stream.NewAggregate()
// 	//agr.Root = NewCustomer("Joe")

// 	event := c.CreateCustomer.New(CustomerCreated{Name: "Joe", Cars: make([]gonvex.AggregateID, 0)})
// 	// if len(a.cars) >= 10 {
// 	// 	return errors.New("not allowed to add more than 10 cars")
// 	// }

// 	return agr.ID, agr.SendEvent(ctx, event)
// 	//return agr.Mutate(ctx, gonvex.NewEvent(PERSON_CREATED, BOUNDED_CONTEXT, person))

// }
