package carpark

import (
	"github.com/alekseev-bro/ddd/pkg/eventstore"
)

type CarRentRejected struct {
	OrderID eventstore.ID[Car]
}

func (ce CarRentRejected) Apply(c *Car) {

}

type CarRented struct {
	OrderID eventstore.ID[Car]
}

func (ce CarRented) Apply(c *Car) {
	c.RentState = NotAvailable
}
