package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
)

type close struct {
}

type closeOrderHandler struct {
	Repo essrv.Root[order.AggregateRoot]
}

func NewCloseOrderHandler(repo essrv.Root[order.AggregateRoot]) *closeOrderHandler {
	return &closeOrderHandler{Repo: repo}
}

func (h *closeOrderHandler) Handle(ctx context.Context, id essrv.ID[order.AggregateRoot], cmd close, idempotencyKey string) error {
	_, err := h.Repo.Execute(ctx, id, func(aggr *order.AggregateRoot) (essrv.Events[order.AggregateRoot], error) {
		return aggr.CloseOrder()
	}, idempotencyKey)
	return err
}
