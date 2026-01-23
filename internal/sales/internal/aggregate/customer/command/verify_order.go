package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type VerifyOrder struct {
	OfCustomer aggregate.ID
	OrderID    aggregate.ID
}

type CustomerUpdater interface {
	Update(ctx context.Context, id aggregate.ID, modify func(state *customer.Customer) (aggregate.Events[customer.Customer], error)) ([]*aggregate.Event[customer.Customer], error)
}

type verifyOrderHandler struct {
	Customers CustomerUpdater
}

func NewVerifyOrderHandler(repo CustomerUpdater) *verifyOrderHandler {
	return &verifyOrderHandler{Customers: repo}
}

func (h *verifyOrderHandler) HandleCommand(ctx context.Context, cmd VerifyOrder) ([]*aggregate.Event[customer.Customer], error) {

	return h.Customers.Update(ctx, cmd.OfCustomer, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.VerifyOrder(cmd.OrderID)
	})
}

type VerifyOrderHandler interface {
	HandleCommand(ctx context.Context, cmd VerifyOrder) ([]*aggregate.Event[customer.Customer], error)
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
