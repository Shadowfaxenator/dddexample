package car

import "github.com/alekseev-bro/ddd/pkg/aggregate"

type Car struct {
	aggregate.Aggregate
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
