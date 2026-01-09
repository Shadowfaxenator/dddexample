package order

import (
	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
)

type Order struct {
	ID         ids.OrderID
	CustomerID ids.CustomerID
	Cars       []OrderLine
	Status     RentOrderStatus
	Deleted    bool
}

func (o *Order) Post() (events.Events[Order], error) {

	return events.New(Posted{
		ID:         o.ID,
		CustomerID: o.CustomerID,
		Cars:       o.Cars,
		Status:     o.Status,
		Deleted:    o.Deleted,
	}), nil

}

func (o *Order) CloseOrder() (events.Events[Order], error) {
	if o.Status != StatusClosed {
		return events.New(Closed{}), nil
	}
	return nil, nil
}
