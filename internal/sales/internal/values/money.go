package values

import (
	"errors"
)

// It's just an artificial money type for demonstration purposes.
type Money struct {
	Currency  string
	Decimal   uint
	Precision uint
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add money with different currencies")
	}
	return Money{
		Currency:  m.Currency,
		Decimal:   m.Decimal + other.Decimal,
		Precision: m.Precision,
	}, nil
}

func NewMoney(currency string, Decimal uint, Precision uint) Money {
	return Money{Currency: currency, Decimal: Decimal, Precision: Precision}
}
