package order

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type OrderLine struct {
	CarID    aggregate.ID // Uses Shared ID
	Price    values.Money // Uses Shared Standard
	Quantity uint         // Primitive
}

func (l OrderLine) Total() values.Money {
	return values.Money{
		Decimal:   l.Price.Decimal * l.Quantity,
		Precision: l.Price.Precision,
		Currency:  l.Price.Currency,
	}
}
