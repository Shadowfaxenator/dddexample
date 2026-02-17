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

type verifyOrderCmdHandler aggregate.CommandHandler[VerifyOrder, customer.Customer]

func NewOrderPostedHandler(handler verifyOrderCmdHandler) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

type orderPostedHandler struct {
	handler verifyOrderCmdHandler
}

func (h *orderPostedHandler) HandleEvent(ctx context.Context, e *order.Posted) error {
	cmd := VerifyOrder{OfCustomer: e.CustomerID, OrderID: e.OrderID}
	_, err := h.handler.HandleCommand(ctx, cmd)

	return err
}
