package order

import (
	"fmt"

	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/internal/sales/internal/values"
)

type OrderLine struct {
	CarID    aggregate.ID // Uses Shared ID
	Price    values.Money // Uses Shared Standard
	Quantity uint         // Primitive
}

type OrderLines []OrderLine

func (l OrderLines) Total() (values.Money, error) {
	var total values.Money
	var err error
	for i, line := range l {
		if i == 0 {
			total = line.Total()
			continue
		}
		fmt.Printf("line.Total(): %v\n", line.Total())
		total, err = total.Add(line.Total())

		if err != nil {
			return values.Money{}, err
		}
	}
	return total, err
}

func (l OrderLine) Total() values.Money {
	return values.Money{
		Decimal:   l.Price.Decimal * l.Quantity,
		Precision: l.Price.Precision,
		Currency:  l.Price.Currency,
	}
}
