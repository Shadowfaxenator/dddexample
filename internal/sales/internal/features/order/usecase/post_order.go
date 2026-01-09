package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
)

type Post struct {
	ID         ids.OrderID
	CustomerID ids.CustomerID
	Cars       []order.OrderLine
}

type postOrderHandler struct {
	Orders events.Store[order.Order]
}

var _ features.CommandHandler[order.Order, Post] = (*postOrderHandler)(nil)

func NewPostOrderHandler(repo events.Store[order.Order]) *postOrderHandler {
	return &postOrderHandler{Orders: repo}
}

func (h *postOrderHandler) Handle(ctx context.Context, id events.ID[order.Order], cmd Post, idempotencyKey string) error {
	_, err := h.Orders.Execute(ctx, id, func(aggr *order.Order) (events.Events[order.Order], error) {

		aggr = &order.Order{
			ID:         cmd.ID,
			CustomerID: cmd.CustomerID,
			Cars:       cmd.Cars,
		}
		return aggr.Post()
	}, idempotencyKey)
	return err
}
