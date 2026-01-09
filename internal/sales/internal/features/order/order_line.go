package order

import (
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/domain/money"
)

type OrderLine struct {
	CarID    ids.CarID   // Uses Shared ID
	Price    money.Money // Uses Shared Standard
	Quantity uint        // Primitive
}

func (l OrderLine) Total() money.Money {
	return money.Money{
		Decimal:   l.Price.Decimal * l.Quantity,
		Precision: l.Price.Precision,
		Currency:  l.Price.Currency,
	}
}
