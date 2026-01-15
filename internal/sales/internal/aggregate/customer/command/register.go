package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type Register struct {
	Customer *customer.Customer
}

func (cmd Register) Execute(c *customer.Customer) (aggregate.Events[customer.Customer], error) {
	return cmd.Customer.Register()
}

type registerHandler struct {
	Customers aggregate.Updater[customer.Customer, *customer.Customer]
}

func NewRegisterHandler(repo aggregate.Updater[customer.Customer, *customer.Customer]) *registerHandler {
	return &registerHandler{Customers: repo}
}

func (h *registerHandler) Handle(ctx context.Context, cmd Register) ([]*aggregate.Event[customer.Customer], error) {

	return h.Customers.Update(ctx, cmd.Customer.ID, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.Register()
	})
}
