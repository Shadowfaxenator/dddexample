package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"

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
	repo essrv.Root[order.AggregateRoot]
}

var _ features.CommandHandler[order.AggregateRoot, Post] = (*postOrderHandler)(nil)

func NewPostOrderHandler(repo essrv.Root[order.AggregateRoot]) *postOrderHandler {
	return &postOrderHandler{repo: repo}
}

func (h *postOrderHandler) Handle(ctx context.Context, id essrv.ID[order.AggregateRoot], cmd Post, idempotencyKey string) error {
	_, err := h.repo.Execute(ctx, id, func(aggr *order.AggregateRoot) (essrv.Events[order.AggregateRoot], error) {

		aggr = &order.AggregateRoot{
			ID:         cmd.ID,
			CustomerID: cmd.CustomerID,
			Cars:       cmd.Cars,
		}
		return aggr.Post()
	}, idempotencyKey)
	return err
}
