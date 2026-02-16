package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	eventstore1 "github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/ddd/pkg/stream"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type Close struct {
	OrderID eventstore1.ID
}

type closeOrderHandler struct {
	Orders orderMutator
}

func NewCloseOrderHandler(repo orderMutator) *closeOrderHandler {
	return &closeOrderHandler{Orders: repo}
}

func (h *closeOrderHandler) HandleCommand(ctx context.Context, cmd Close) ([]stream.MsgMetadata, error) {

	return h.Orders.Mutate(ctx, cmd.OrderID, func(state *order.Order) (eventstore1.Events[order.Order], error) {
		return state.Close()
	})
}

type orderRejectedHandler struct {
	CloseOrderHandler aggregate.CommandHandler[order.Order, Close]
}

func NewOrderRejectedHandler(h aggregate.CommandHandler[order.Order, Close]) *orderRejectedHandler {

	return &orderRejectedHandler{CloseOrderHandler: h}
}

func (h *orderRejectedHandler) HandleEvent(ctx context.Context, e *customer.OrderRejected) error {
	cmd := Close{OrderID: e.OrderID}
	_, err := h.CloseOrderHandler.HandleCommand(ctx, cmd)

	return err
}
