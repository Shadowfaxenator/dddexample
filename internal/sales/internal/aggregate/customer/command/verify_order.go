package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type VerifyOrder struct {
	OfCustomer aggregate.ID
	OrderID    aggregate.ID
}

type verifyOrderHandler struct {
	Customers customerMutator
}

func NewVerifyOrderHandler(repo customerMutator) *verifyOrderHandler {
	return &verifyOrderHandler{Customers: repo}
}

func (h *verifyOrderHandler) HandleCommand(ctx context.Context, cmd VerifyOrder) ([]stream.EventMetadata, error) {
	return h.Customers.Mutate(ctx, cmd.OfCustomer, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.VerifyOrder(cmd.OrderID)
	})
}

type VerifyOrderHandler interface {
	HandleCommand(ctx context.Context, cmd VerifyOrder) ([]stream.EventMetadata, error)
}

func NewOrderPostedHandler(handler VerifyOrderHandler) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

type orderPostedHandler struct {
	handler VerifyOrderHandler
}

func (h *orderPostedHandler) HandleEvent(ctx context.Context, e *order.Posted) error {
	cmd := VerifyOrder{OfCustomer: e.CustomerID, OrderID: e.OrderID}
	_, err := h.handler.HandleCommand(ctx, cmd)

	return err
}
