package orders

import (
	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain"
)

type Order struct {
	ID         domain.OrderID
	CustomerID domain.CustomerID
	Cars       map[domain.CarID]struct{}
	Status     RentOrderStatus
	OrderLines []OrderLine
	Deleted    bool
}

func (o *Order) PostOrder(ord *Order) essrv.Events[Order] {
	if o.ID.IsZero() {
		return essrv.NewEvents(OrderPosted{Order: ord})
	}
	return nil
}

func (o *Order) CloseOrder() essrv.Events[Order] {
	if o.Status != Closed {
		return essrv.NewEvents(OrderClosed{})
	}
	return nil
}
