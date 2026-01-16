package command

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type CustomerRegisterHandler aggregate.CommandHandler[customer.Customer, Register]
