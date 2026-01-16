package command

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/order"
)

type OrderpostHandler aggregate.CommandHandler[order.Order, Post]
