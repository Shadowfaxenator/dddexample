package cars

import "github.com/alekseev-bro/dddexample/internal/sales/internal/domain"

type Car struct {
	ID           domain.CarID
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
