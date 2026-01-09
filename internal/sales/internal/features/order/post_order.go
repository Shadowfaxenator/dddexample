package order

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
)

type Post struct {
	*orders.Order
}

type postOrderHandler struct {
	repo essrv.Root[orders.Order]
}

var _ features.CommandHandler[orders.Order, Post] = (*postOrderHandler)(nil)

func NewPostOrderHandler(repo essrv.Root[orders.Order]) *postOrderHandler {
	return &postOrderHandler{repo: repo}
}

func (h *postOrderHandler) Handle(ctx context.Context, id essrv.ID[orders.Order], cmd Post, idempotencyKey string) error {
	_, err := h.repo.Execute(ctx, id, func(aggr *orders.Order) (essrv.Events[orders.Order], error) {
		return aggr.PostOrder(cmd.Order), nil
	}, id.String())
	return err
}
