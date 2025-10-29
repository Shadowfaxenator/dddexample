package sales

import (
	"context"
	"ddd/pkg/aggregate"
)

type CustomerService struct {
	customer Aggregate[Customer]
	order    Aggregate[Order]
}

func NewCustomerService(ctx context.Context, customer Aggregate[Customer], order Aggregate[Order]) *CustomerService {
	s := &CustomerService{
		customer: customer,
		order:    order,
	}
	s.order.Subscribe(ctx, "sales_customer_service", func(e aggregate.Event[Order]) error {
		switch ev := e.(type) {
		case *OrderCreated:
			return s.customer.Command(ctx, ev.Order.CustomerID, ValidateOrder{OrderID: ev.Order.ID})
		}

		return nil
	}, false)
	return s
}
