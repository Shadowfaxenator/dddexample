package car

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/contracts/v1/carpark"
)

type Arrived struct {
	Car *Car
}

func (ca *Arrived) ToArrivedV1() *carpark.CarArrived {
	return &carpark.CarArrived{
		Car: ca.Car.ToCarV1(),
	}
}

func (ca *Arrived) Evolve(c *Car) {
	*c = *ca.Car
}

type CarRentRejected struct {
	OrderID aggregate.ID
}

func (ce *CarRentRejected) Evolve(c *Car) {

}

type CarRented struct {
	OrderID aggregate.ID
}

func (ce *CarRented) Evolve(c *Car) {
	c.RentState = NotAvailable
}
