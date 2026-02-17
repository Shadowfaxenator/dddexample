package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type Close struct {
	OrderID aggregate.ID
}

type closeOrderHandler struct {
	Orders orderMutator
}

func NewCloseOrderHandler(repo orderMutator) *closeOrderHandler {
	return &closeOrderHandler{Orders: repo}
}

func (h *closeOrderHandler) HandleCommand(ctx context.Context, cmd Close) ([]stream.EventMetadata, error) {

	return h.Orders.Mutate(ctx, cmd.OrderID, func(state *order.Order) (aggregate.Events[order.Order], error) {
		return state.Close()
	})
}

type closeOrderCmdHandler aggregate.CommandHandler[Close, order.Order]

type orderRejectedHandler struct {
	CloseOrderHandler closeOrderCmdHandler
}

func NewOrderRejectedHandler(h closeOrderCmdHandler) *orderRejectedHandler {

	return &orderRejectedHandler{CloseOrderHandler: h}
}

func (h *orderRejectedHandler) HandleEvent(ctx context.Context, e *customer.OrderRejected) error {
	cmd := Close{OrderID: e.OrderID}
	_, err := h.CloseOrderHandler.HandleCommand(ctx, cmd)

	return err
}
