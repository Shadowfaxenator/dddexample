package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
)

type verifyOrder struct {
	OrderID ids.OrderID
}

type verifyOrderHandler struct {
	Repo essrv.Root[customer.AggregateRoot]
}

func NewVerifyOrderHandler(repo essrv.Root[customer.AggregateRoot]) *verifyOrderHandler {
	return &verifyOrderHandler{Repo: repo}
}

func (h *verifyOrderHandler) Handle(ctx context.Context, id essrv.ID[customer.AggregateRoot], cmd verifyOrder, idempotencyKey string) error {
	_, err := h.Repo.Execute(ctx, id, func(c *customer.AggregateRoot) (essrv.Events[customer.AggregateRoot], error) {
		return c.VerifyOrder(cmd.OrderID)
	}, idempotencyKey)
	return err
}

func NewOrderPostedHandler(handler features.CommandHandler[customer.AggregateRoot, verifyOrder]) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

type orderPostedHandler struct {
	handler features.CommandHandler[customer.AggregateRoot, verifyOrder]
}

func (h *orderPostedHandler) Handle(ctx context.Context, eventID string, event order.Posted) error {
	return h.handler.Handle(
		ctx,
		essrv.ID[customer.AggregateRoot](event.CustomerID),
		verifyOrder{OrderID: event.ID},
		eventID,
	)
}
