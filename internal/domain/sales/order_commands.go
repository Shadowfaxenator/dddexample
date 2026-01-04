package sales

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

var ErrOrderAlreadyClosed = errors.New("order already closed")

func (o *Order) Close() (eventstore.Events[Order], error) {
	if o.Status != Closed {
		return eventstore.NewEvents(OrderClosed{}), nil
	}
	return nil, ErrOrderAlreadyClosed
}

func (o *Order) Create(customerID eventstore.ID[Customer]) (eventstore.Events[Order], error) {
	o.CustomerID = customerID
	o.Cars = map[eventstore.ID[Car]]struct{}{}
	return eventstore.NewEvents(OrderCreated{Order: *o}), nil
}
