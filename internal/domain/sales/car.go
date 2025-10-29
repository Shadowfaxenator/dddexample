package sales

import "ddd/pkg/aggregate"

type Car struct {
	aggregate.ID[Car]
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
