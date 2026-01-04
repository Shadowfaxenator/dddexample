package sales

import "github.com/alekseev-bro/ddd/pkg/eventstore"

type Car struct {
	eventstore.ID[Car]
	Make         string
	Model        string
	Year         uint
	Price        uint
	Availability bool
}
