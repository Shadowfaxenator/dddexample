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

func (cmd Close) Execute(o *order.Order) (aggregate.Events[order.Order], error) {
	return o.Close()
}

type closeOrderHandler struct {
	Orders aggregate.Updater[order.Order, *order.Order]
}

func NewCloseOrderHandler(repo aggregate.Updater[order.Order, *order.Order]) *closeOrderHandler {
	return &closeOrderHandler{Orders: repo}
}

func (h *closeOrderHandler) Handle(ctx context.Context, cmd Close) ([]*aggregate.Event[order.Order], error) {

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

func (h *orderRejectedHandler) HandleEvent(ctx context.Context, e customer.OrderRejected) error {
	cmd := Close{OrderID: e.OrderID}
	_, err := h.CloseOrderHandler.Handle(ctx, cmd)
	return err
}
