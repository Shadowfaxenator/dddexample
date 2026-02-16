package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type Post struct {
	Order *order.Order
}

type postOrderHandler struct {
	Orders orderMutator
}

func NewPostOrderHandler(repo orderMutator) *postOrderHandler {
	return &postOrderHandler{Orders: repo}
}

func (h *postOrderHandler) HandleCommand(ctx context.Context, cmd Post) ([]stream.EventMetadata, error) {

	return h.Orders.Mutate(ctx, cmd.Order.ID, func(state *order.Order) (aggregate.Events[order.Order], error) {
		return state.Post(cmd.Order)
	})
}
