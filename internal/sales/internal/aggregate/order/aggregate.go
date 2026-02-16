package order

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	eventstore1 "github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type Order struct {
	ID         eventstore1.ID
	CustomerID eventstore1.ID
	Cars       OrderLines
	Total      values.Money
	Status     RentOrderStatus
}

func New(customerID eventstore1.ID, cars OrderLines) *Order {
	total, err := cars.Total()
	if err != nil {
		panic(err)
	}
	id, err := eventstore1.NewID()
	if err != nil {
		panic(err)
	}

	o := &Order{
		ID:         id,
		CustomerID: customerID,
		Cars:       cars,
		Total:      total,
	}
	return o
}

func (o *Order) Post(ord *Order) (eventstore1.Events[Order], error) {
	if o.Status != StatusNew {
		return nil, aggregate.ErrAggregateAlreadyExists
	}
	return eventstore1.NewEvents(&Posted{
		OrderID:    ord.ID,
		CustomerID: ord.CustomerID,
		Cars:       ord.Cars,
		Status:     ord.Status,
		Total:      ord.Total,
	}), nil

}

func (o *Order) Close() (eventstore1.Events[Order], error) {
	if o.Status != StatusClosed {
		return eventstore1.NewEvents(&Closed{}), nil
	}
	return nil, nil
}
