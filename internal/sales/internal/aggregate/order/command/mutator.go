package command

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type orderMutator aggregate.Mutator[order.Order, *order.Order]
