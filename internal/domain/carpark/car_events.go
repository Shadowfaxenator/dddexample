package carpark

import (
	"ddd/pkg/aggregate"
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

type CarCreated struct {
	Car Car
}

func (cc CarCreated) Apply(c *Car) {
	*c = cc.Car
}

func (CarCreated) Type() string {
	return "CarCreated"
}
