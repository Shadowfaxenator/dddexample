package carpark

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
)

type CarRentRejected struct {
	OrderID aggregate.ID[Car]
}

func (ce CarRentRejected) Apply(c *Car) {

}

type CarRented struct {
	OrderID aggregate.ID[Car]
}

func (ce CarRented) Apply(c *Car) {
	c.RentState = NotAvailable
}
