package sales

import (
	"context"
	"ddd/pkg/aggregate"
	"errors"
	"fmt"
	"log/slog"
)

func (s *SubDomain) CustomerService(ctx context.Context) {
	s.order.Subscribe(ctx, "sales_customer_service", func(e aggregate.Applyer[Order]) error {
		switch ev := e.(type) {
		case *OrderCreated:
			err := s.ValidateOrder(ctx, ev.Order.CustomerID, ev.Order.ID)
			if err != nil {
				switch {
				case errors.Is(err, ErrValidateAge):
					err := s.CloseOrder(ctx, ev.Order.ID)
					if err != nil {
						return fmt.Errorf("customer service: %w", err)
					}
					return nil
				case errors.Is(err, ErrMaxOrders):
					err := s.CloseOrder(ctx, ev.Order.ID)
					if err != nil {
						return fmt.Errorf("customer service: %w", err)
					}
					return nil
				}
				return fmt.Errorf("customer service: %w", err)
			}
		case *OrderClosed:
			err := s.removeFromOneOrder(ctx, ev.CustID)
			if err != nil {
				var e *ValidateOrdersError
				switch {
				case errors.As(err, &e):
					slog.Warn(e.Error())
					return nil
				default:
					return fmt.Errorf("customer service: %w", err)
				}
			}
		}
		return nil
	}, false)
}
