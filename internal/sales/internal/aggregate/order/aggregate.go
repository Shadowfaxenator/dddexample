package order

import (
	"errors"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type Order struct {
	aggregate.Aggregate
	CustomerID aggregate.ID
	Cars       []OrderLine
	Status     RentOrderStatus
}

func New(id, customerID aggregate.ID, cars []OrderLine) *Order {
	o := &Order{
		CustomerID: customerID,
		Cars:       cars,
	}
	o.ID = id
	return o
}

func (o *Order) Post(ord *Order) (aggregate.Events[Order], error) {
	if o.Exists {
		return nil, errors.New("order exists")
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
