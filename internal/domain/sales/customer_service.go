package sales

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type CustomerService struct {
	Order eventstore.EventStore[Order]
}

func (c CustomerService) Handle(ctx context.Context, eventID eventstore.EventID[Customer], e eventstore.Event[Customer]) error {
	switch ev := e.(type) {
	case *OrderAccepted:
		_, err := c.Order.Update(ctx, ev.OrderID, eventID.String(), func(o *Order) (eventstore.Events[Order], error) {
			return o.Close()
		})
		return err

	}
	return nil
}
