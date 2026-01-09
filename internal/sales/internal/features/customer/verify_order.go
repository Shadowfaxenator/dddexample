package customer

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/customers"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/orders"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/features"
)

type verifyOrder struct {
	OrderID domain.OrderID
}

type verifyOrderHandler struct {
	Repo essrv.Root[customers.Customer]
}

func NewVerifyOrderHandler(repo essrv.Root[customers.Customer]) *verifyOrderHandler {
	return &verifyOrderHandler{Repo: repo}
}

func (h *verifyOrderHandler) Handle(ctx context.Context, id essrv.ID[customers.Customer], cmd verifyOrder, idempotencyKey string) error {
	_, err := h.Repo.Execute(ctx, id, func(c *customers.Customer) (essrv.Events[customers.Customer], error) {
		return c.VerifyOrder(cmd.OrderID), nil
	}, idempotencyKey)
	return err
}

type orderPostedHandler struct {
	handler features.CommandHandler[customers.Customer, verifyOrder]
}

func NewOrderPostedHandler(handler features.CommandHandler[customers.Customer, verifyOrder]) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

func (h *orderPostedHandler) Handle(ctx context.Context, eventID string, event orders.OrderPosted) error {
	return h.handler.Handle(
		ctx,
		essrv.ID[customers.Customer](event.CustomerID),
		verifyOrder{OrderID: event.ID},
		eventID,
	)
}
