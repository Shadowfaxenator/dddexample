package carpark

import (
	"github.com/alekseev-bro/ddd/pkg/events"
)

type CarRentRejected struct {
	OrderID events.ID[Car]
}

func (ce CarRentRejected) Apply(c *Car) {

}

type CarRented struct {
	OrderID events.ID[Car]
}

func (ce CarRented) Apply(c *Car) {
	c.RentState = NotAvailable
}
