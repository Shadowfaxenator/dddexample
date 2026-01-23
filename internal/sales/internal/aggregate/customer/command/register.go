package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type Register struct {
	Customer *customer.Customer
}

type registerHandler struct {
	Customers RegisterHandler
}

type RegisterHandler interface {
	Update(ctx context.Context, id aggregate.ID, modify func(state *customer.Customer) (aggregate.Events[customer.Customer], error)) ([]*aggregate.Event[customer.Customer], error)
}

func NewRegisterHandler(repo RegisterHandler) *registerHandler {
	return &registerHandler{Customers: repo}
}

func (h *registerHandler) HandleCommand(ctx context.Context, cmd Register) ([]*aggregate.Event[customer.Customer], error) {

	return h.Customers.Update(ctx, cmd.Customer.ID, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.Register(cmd.Customer)
	})
}
