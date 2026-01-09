package order

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/orders"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
)

type close struct {
}

var _ features.CommandHandler[orders.Order, close] = (*closeOrderHandler)(nil)

type closeOrderHandler struct {
	Repo essrv.Root[orders.Order]
}

func NewCloseOrderHandler(repo essrv.Root[orders.Order]) *closeOrderHandler {
	return &closeOrderHandler{Repo: repo}
}

func (h *closeOrderHandler) Handle(ctx context.Context, id essrv.ID[orders.Order], cmd close, idempotencyKey string) error {
	_, err := h.Repo.Execute(ctx, id, func(aggr *orders.Order) (essrv.Events[orders.Order], error) {
		return aggr.CloseOrder(), nil
	}, idempotencyKey)
	return err
}
