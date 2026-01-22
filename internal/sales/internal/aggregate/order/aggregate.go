package order

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Order struct {
	ID         aggregate.ID
	CustomerID aggregate.ID
	Cars       []OrderLine
	Status     RentOrderStatus
}

func New(customerID aggregate.ID, cars []OrderLine) *Order {
	o := &Order{
		ID:         aggregate.NewID(),
		CustomerID: customerID,
		Cars:       cars,
	}
	return o
}

func (o *Order) Post(ord *Order) (aggregate.Events[Order], error) {
	if o.Status != StatusNew {
		return nil, aggregate.ErrAggregateAlreadyExists
	}
	return aggregate.NewEvents(&Posted{
		OrderID:    ord.ID,
		CustomerID: ord.CustomerID,
		Cars:       ord.Cars,
		Status:     ord.Status,
	}), nil

}

func (o *Order) Close() (aggregate.Events[Order], error) {
	if o.Status != StatusClosed {
		return aggregate.NewEvents(&Closed{}), nil
	}
	return nil, nil
}
