package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/domain"
)

type CustomerService struct {
	Order domain.Aggregate[Order]
}

func (c *CustomerService) Handle(ctx context.Context, eventID domain.EventID[Customer], e domain.Event[Customer]) error {
	switch ev := e.(type) {
	case *OrderAccepted:
		_, err := c.Order.Execute(ctx, string(eventID), &CloseOrder{OrderID: ev.OrderID})
		return err

	}
	return nil
}
