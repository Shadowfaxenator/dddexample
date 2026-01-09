package usecase

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
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
	repo essrv.Root[customer.AggregateRoot]
}

func NewRegisterHandler(repo essrv.Root[customer.AggregateRoot]) *registerHandler {
	return &registerHandler{repo: repo}
}

func (h *registerHandler) Handle(ctx context.Context, id essrv.ID[customer.AggregateRoot], cmd Register, idempotencyKey string) error {
	_, err := h.repo.Execute(ctx, id, func(aggr *customer.AggregateRoot) (essrv.Events[customer.AggregateRoot], error) {

		aggr = &customer.AggregateRoot{
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
