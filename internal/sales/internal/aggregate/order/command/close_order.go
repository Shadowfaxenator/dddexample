package command

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type Close struct {
	OrderID aggregate.ID
}

type closeOrderHandler struct {
	Orders OrderUpdater
}

type OrderUpdater interface {
	Update(ctx context.Context, id aggregate.ID, modify func(state *order.Order) (aggregate.Events[order.Order], error)) ([]*aggregate.Event[order.Order], error)
}

func NewCloseOrderHandler(repo OrderUpdater) *closeOrderHandler {
	return &closeOrderHandler{Orders: repo}
}

func (h *closeOrderHandler) HandleCommand(ctx context.Context, cmd Close) ([]*aggregate.Event[order.Order], error) {

	return h.Orders.Update(ctx, cmd.OrderID, func(state *order.Order) (aggregate.Events[order.Order], error) {
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
