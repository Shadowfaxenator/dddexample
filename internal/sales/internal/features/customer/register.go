package customer

import (
	"context"

	"github.com/alekseev-bro/ddd/pkg/essrv"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/customers"
)

type Register struct {
	*customers.Customer
}

type RegisterHandler struct {
	repo essrv.Root[customers.Customer]
}

func NewRegisterHandler(repo essrv.Root[customers.Customer]) *RegisterHandler {
	return &RegisterHandler{repo: repo}
}

func (h *RegisterHandler) Handle(ctx context.Context, id essrv.ID[customers.Customer], cmd Register, idempotencyKey string) error {
	_, err := h.repo.Execute(ctx, id, func(aggr *customers.Customer) (essrv.Events[customers.Customer], error) {
		return aggr.Register(cmd.Customer), nil
	}, id.String())
	return err
}
