package order

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type Order struct {
	ID         aggregate.ID
	CustomerID aggregate.ID
	Cars       OrderLines
	Total      values.Money
	Status     RentOrderStatus
}

func New(customerID aggregate.ID, cars OrderLines) *Order {
	total, err := cars.Total()
	if err != nil {
		panic(err)
	}
	id, err := aggregate.NewID()
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

func (o *Order) Post(ord *Order) (aggregate.Events[Order], error) {
	if o.Status != StatusNew {
		return nil, aggregate.ErrAlreadyExists
	}
	return aggregate.NewEvents(&Posted{
		OrderID:    ord.ID,
		CustomerID: ord.CustomerID,
		Cars:       ord.Cars,
		Status:     ord.Status,
		Total:      ord.Total,
	}), nil

}

func (o *Order) Close() (aggregate.Events[Order], error) {
	if o.Status != StatusClosed {
		return aggregate.NewEvents(&Closed{}), nil
	}
	return nil, nil
}
