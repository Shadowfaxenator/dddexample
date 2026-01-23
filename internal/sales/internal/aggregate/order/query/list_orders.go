package query

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer/query"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type Order struct {
	ID       aggregate.ID
	UserID   aggregate.ID
	Total    values.Money
	UserName string
}

type AllLister interface {
	ListAll() ([]Order, error)
}

type CustomerGetter interface {
	GetCustomer(id aggregate.ID) (*query.Customer, bool)
}

type OrderListProjector struct {
	Orders    *MemOrders
	Customers CustomerGetter
}

func NewOrderListProjector(customers CustomerGetter, orders *MemOrders) *OrderListProjector {
	return &OrderListProjector{
		Orders:    orders,
		Customers: customers,
	}
}

type OrderListQueryHandler struct {
	Orders AllLister
}

type MemOrders struct {
	Orders []Order
}

func (m *MemOrders) ListAll() ([]Order, error) {
	return m.Orders, nil
}

func (m *MemOrders) AddOrder(order Order) {

	m.Orders = append(m.Orders, order)
}

func NewMemOrders() *MemOrders {
	return &MemOrders{
		Orders: make([]Order, 0),
	}
}

func (h *OrderListProjector) HandleEvents(ctx context.Context, event aggregate.Evolver[order.Order]) error {
	switch ev := event.(type) {
	case *order.Posted:
		if cust, ok := h.Customers.GetCustomer(ev.CustomerID); ok {
			h.Orders.AddOrder(Order{
				ID:       ev.OrderID,
				Total:    ev.Total,
				UserID:   ev.CustomerID,
				UserName: cust.Name,
			})
		} else {
			h.Orders.AddOrder(Order{
				ID:       ev.OrderID,
				Total:    ev.Total,
				UserID:   ev.CustomerID,
				UserName: "NaN",
			})
		}

	}
	return nil
}
