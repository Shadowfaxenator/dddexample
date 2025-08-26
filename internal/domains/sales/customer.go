package sales

import (
	"context"
	"ttt/gonvex"
	"ttt/gonvex/schema"
)

func NewCustomerAggregate(ctx context.Context) *CustomerRoot {
	return &CustomerRoot{
		gonvex.NewAggregateRoot[Customer](ctx),
	}
}

type CustomerRoot struct {
	gonvex.AggregateRoot[*Customer]
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
	Dogs   []string
}

type Customer struct {
	Name     schema.Field[string]
	Adresses schema.List[Address]
	Cars     schema.List[gonvex.ID]
}

func (p *Customer) Events() []gonvex.Event {
	return []gonvex.Event{
		CustomerCarAdded{},
		CustomerCreated{},
	}
}

func (p *Customer) Reduce(event gonvex.Event) error {

	switch ev := event.(type) {
	case *CustomerCarAdded:
		p.Cars.Push(ev.CarID)
	case *CustomerCreated:

		per := ev.Customer

		*p = per
	}

	return nil
}

func (c *CustomerRoot) Create(ctx context.Context, name string) (gonvex.ID, error) {
	return c.Command(ctx, nil, func(c *Customer) (gonvex.Event, error) {

		return &CustomerCreated{Customer: Customer{Name: schema.NewField(name), Cars: schema.NewList[gonvex.ID]()}, AggID: gonvex.NewID()}, nil

	})
}

func (c *CustomerRoot) AddCar(ctx context.Context, id gonvex.ID) error {
	_, err := c.Command(ctx, &id, func(c *Customer) (gonvex.Event, error) {

		// if len(c.Cars) >= 10 {
		// 	return nil, fmt.Errorf("can't add cars >= 10")
		// }
		return &CustomerCarAdded{AggID: id, CarID: gonvex.NewID()}, nil
	})

	return err

}

//agr.Root = NewCustomer("Joe")

// if len(a.cars) >= 10 {
// 	return errors.New("not allowed to add more than 10 cars")
// }

//return agr.Mutate(ctx, gonvex.NewEvent(PERSON_CREATED, BOUNDED_CONTEXT, person))
