package sales

import (
	"context"
	"ddd/pkg/aggregate"
)

type OrderService struct {
	customer Aggregate[Customer]
	order    Aggregate[Order]
}

func NewOrderService(ctx context.Context, customer Aggregate[Customer], order Aggregate[Order]) *OrderService {
	s := &OrderService{
		customer: customer,
		order:    order,
	}
	s.customer.Subscribe(ctx, "sales_order_service", func(e aggregate.Event[Customer]) error {
		switch ev := e.(type) {
		case *OrderAccepted:
			return s.order.Command(ctx, ev.OrderID, CloseOrder{OrderID: ev.OrderID})

		}
		return nil
	}, false)
	return s
}
