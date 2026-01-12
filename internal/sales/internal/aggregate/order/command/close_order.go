package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type close struct {
}

type closeOrderHandler struct {
	Orders events.Executer[order.Order]
}

func NewCloseOrderHandler(repo events.Executer[order.Order]) *closeOrderHandler {
	return &closeOrderHandler{Orders: repo}
}

func (h *closeOrderHandler) Handle(ctx context.Context, id events.ID[order.Order], cmd close, idempotencyKey string) error {
	_, err := h.Orders.Execute(ctx, id, func(aggr *order.Order) (events.Events[order.Order], error) {
		return aggr.CloseOrder()
	}, idempotencyKey)
	return err
}

type orderRejectedHandler struct {
	Repo events.Store[order.Order]
}

func NewOrderRejectedHandler(repo events.Store[order.Order]) *orderRejectedHandler {

	return &orderRejectedHandler{Repo: repo}
}

func (h *orderRejectedHandler) Handle(ctx context.Context, eventID string, e customer.OrderRejected) error {
	_, err := h.Repo.Execute(ctx, events.ID[order.Order](e.OrderID), func(aggr *order.Order) (events.Events[order.Order], error) {
		return aggr.CloseOrder()
	}, eventID)
	return err
}
