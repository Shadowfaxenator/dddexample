package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/events"

	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/features/customer"
)

type Register struct {
	ID           ids.CustomerID
	Name         string
	Age          uint
	Addresses    []customer.Address
	ActiveOrders uint
}

type registerHandler struct {
	Customers events.Store[customer.Customer]
}

func NewRegisterHandler(repo events.Store[customer.Customer]) *registerHandler {
	return &registerHandler{Customers: repo}
}

func (h *registerHandler) Handle(ctx context.Context, id events.ID[customer.Customer], cmd Register, idempotencyKey string) error {
	_, err := h.Customers.Execute(ctx, id, func(aggr *customer.Customer) (events.Events[customer.Customer], error) {

		aggr = &customer.Customer{
			ID:           cmd.ID,
			Name:         cmd.Name,
			Age:          cmd.Age,
			Addresses:    cmd.Addresses,
			ActiveOrders: cmd.ActiveOrders,
		}
		return aggr.Register()

	}, idempotencyKey)
	return err
}
