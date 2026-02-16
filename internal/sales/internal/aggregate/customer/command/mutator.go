package command

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/aggregate/customer"
)

type customerMutator aggregate.Mutator[customer.Customer, *customer.Customer]
