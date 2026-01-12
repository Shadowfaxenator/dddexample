package query

import (
	"context"
	"time"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type OrderListProjection struct {
	ID        values.OrderID
	Total     values.Money
	CreatedAt time.Time
	UserName  string
}

type AllLister interface {
	ListAll() ([]OrderListProjection, error)
}

type OrderListProjector struct {
	Orders *MemOrders
}

type OrderListQueryHandler struct {
	Orders AllLister
}

type MemOrders struct {
	Orders []OrderListProjection
}

func (m *MemOrders) ListAll() ([]OrderListProjection, error) {
	return m.Orders, nil
}

func NewMemOrders() *MemOrders {
	return &MemOrders{
		Orders: make([]OrderListProjection, 0),
	}
}

func (h *OrderListProjector) Handle(ctx context.Context, eventID string, event any) error {
	switch event.(type) {
	case order.Posted:

	}
	return nil
}
