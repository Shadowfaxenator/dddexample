package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type Register struct {
	Customer *customer.Customer
}

type registerHandler struct {
	Customers customerMutator
}

func NewRegisterHandler(repo customerMutator) *registerHandler {
	return &registerHandler{Customers: repo}
}

func (h *registerHandler) HandleCommand(ctx context.Context, cmd Register) ([]stream.EventMetadata, error) {

	return h.Customers.Mutate(ctx, cmd.Customer.ID, func(state *customer.Customer) (aggregate.Events[customer.Customer], error) {
		return state.Register(cmd.Customer)
	})
}
