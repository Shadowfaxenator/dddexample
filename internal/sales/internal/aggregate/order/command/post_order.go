package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type Post struct {
	Order *order.Order
}

type postOrderHandler struct {
	Orders OrderUpdater
}

func NewPostOrderHandler(repo OrderUpdater) *postOrderHandler {
	return &postOrderHandler{Orders: repo}
}

func (h *postOrderHandler) HandleCommand(ctx context.Context, cmd Post) ([]*aggregate.Event[order.Order], error) {

	return h.Orders.Update(ctx, cmd.Order.ID, func(state *order.Order) (aggregate.Events[order.Order], error) {
		return state.Post(cmd.Order)
	})
}
