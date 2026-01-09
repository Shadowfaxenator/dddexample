package orders

import "github.com/alekseev-bro/dddexample/internal/sales/internal/domain"

type OrderLine struct {
	ProductID domain.ProductID // Uses Shared ID
	Price     domain.Money     // Uses Shared Standard
	Quantity  uint             // Primitive
}

func (l OrderLine) Total() domain.Money {
	return domain.Money{
		Decimal:   l.Price.Decimal * l.Quantity,
		Precision: l.Price.Precision,
		Currency:  l.Price.Currency,
	}
}
