package car

import "github.com/alekseev-bro/dddexample/internal/sales/internal/values"

type Car struct {
	ID           values.CarID
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
