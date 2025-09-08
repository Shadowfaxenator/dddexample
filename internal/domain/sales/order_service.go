package sales

import (
	"context"
	"ddd/pkg/aggregate"
	"fmt"
)

func (s *SubDomain) OrderService(ctx context.Context) {
	s.customer.Subscribe(ctx, "sales_order_service", func(e aggregate.Applyer[Customer]) error {
		switch ev := e.(type) {
		case *OrderAccepted:
			err := s.MarkVerrified(ctx, ev.OrderID)
			if err != nil {
				return fmt.Errorf("customer service: %w", err)
			}
		}
		return nil
	}, false)
}
