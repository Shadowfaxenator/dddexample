package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/order"
)

type verifyOrder struct {
	OrderID ids.OrderID
}

type verifyOrderHandler struct {
	Customers events.Store[customer.Customer]
}

func NewVerifyOrderHandler(repo events.Store[customer.Customer]) *verifyOrderHandler {
	return &verifyOrderHandler{Customers: repo}
}

func (h *verifyOrderHandler) Handle(ctx context.Context, id events.ID[customer.Customer], cmd verifyOrder, idempotencyKey string) error {
	_, err := h.Customers.Execute(ctx, id, func(c *customer.Customer) (events.Events[customer.Customer], error) {
		return c.VerifyOrder(cmd.OrderID)
	}, idempotencyKey)
	return err
}

func NewOrderPostedHandler(handler features.CommandHandler[customer.Customer, verifyOrder]) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

type orderPostedHandler struct {
	handler features.CommandHandler[customer.Customer, verifyOrder]
}

func (h *orderPostedHandler) Handle(ctx context.Context, eventID string, e order.Posted) error {
	return h.handler.Handle(
		ctx,
		events.ID[customer.Customer](e.CustomerID),
		verifyOrder{OrderID: e.ID},
		eventID,
	)
}
