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

func (cmd VerifyOrder) Execute(c *customer.Customer) (aggregate.Events[customer.Customer], error) {
	return c.VerifyOrder(cmd.OrderID)
}

type verifyOrderHandler struct {
	Customers aggregate.Updater[customer.Customer, *customer.Customer]
}

func NewVerifyOrderHandler(repo aggregate.Updater[customer.Customer, *customer.Customer]) *verifyOrderHandler {
	return &verifyOrderHandler{Customers: repo}
}

func (h *verifyOrderHandler) Handle(ctx context.Context, cmd VerifyOrder) ([]*aggregate.Event[customer.Customer], error) {
	return h.Customers.Update(ctx, cmd.OfCustomer, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.VerifyOrder(cmd.OrderID)
	})
}

func NewOrderPostedHandler(handler aggregate.CommandHandler[customer.Customer, VerifyOrder]) *orderPostedHandler {
	return &orderPostedHandler{handler: handler}
}

type orderPostedHandler struct {
	handler aggregate.CommandHandler[customer.Customer, VerifyOrder]
}

func (h *orderPostedHandler) HandleEvent(ctx context.Context, e order.Posted) error {
	cmd := VerifyOrder{OfCustomer: e.CustomerID, OrderID: e.OrderID}
	_, err := h.handler.Handle(ctx, cmd)
	return err
}
