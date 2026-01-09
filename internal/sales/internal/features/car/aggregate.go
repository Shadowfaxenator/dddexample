package car

import "github.com/alekseev-bro/dddexample/internal/sales/internal/domain/ids"

type Car struct {
	ID           ids.CarID
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
