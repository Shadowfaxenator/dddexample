package sales

import (
	"context"
	"ddd/pkg/aggregate"
)

type SubDomain struct {
	customer *aggregate.Aggregate[Customer]
	order    *aggregate.Aggregate[Order]
}

func NewSubDomain(ctx context.Context) *SubDomain {
	customer := aggregate.New[Customer](ctx)
	aggregate.RegisterEvent[CustomerCreated](customer)
	aggregate.RegisterEvent[OrderAccepted](customer)
	//aggregate.RegisterCommand[CreateCustomer](customer)
	order := aggregate.New[Order](ctx)
	aggregate.RegisterEvent[OrderCreated](order)
	aggregate.RegisterEvent[OrderClosed](order)

	aggregate.RegisterEvent[OrderVerified](order)
	//aggregate.RegisterCommand[CreateOrder](order)
	c := &SubDomain{
		customer: customer,
		order:    order,
	}
	c.CustomerService(ctx)
	c.OrderService(ctx)
	return c
}
