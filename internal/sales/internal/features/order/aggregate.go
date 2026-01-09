package order

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
)

type AggregateRoot struct {
	ID         ids.OrderID
	CustomerID ids.CustomerID
	Cars       []OrderLine
	Status     RentOrderStatus
	Deleted    bool
}

func (o *AggregateRoot) Post() (essrv.Events[AggregateRoot], error) {

	return essrv.NewEvents(Posted{
		ID:         o.ID,
		CustomerID: o.CustomerID,
		Cars:       o.Cars,
		Status:     o.Status,
		Deleted:    o.Deleted,
	}), nil

}

func (o *AggregateRoot) CloseOrder() (essrv.Events[AggregateRoot], error) {
	if o.Status != StatusClosed {
		return essrv.NewEvents(Closed{}), nil
	}
	return nil, nil
}
