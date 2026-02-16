package car

import (
	"github.com/alekseev-bro/ddd/pkg/aggregate"
	"github.com/alekseev-bro/dddexample/contracts/v1/carpark"
	"github.com/google/uuid"
)

// import (
// 	"encoding/json"
// 	"fmt"
// 	"ttt/gonvex"

// 	"github.com/google/uuid"
// )

//	func newCar(model string, brand string) *Car {
//		return &Car{ID: gonvex.NewAggregateID(), VIN: uuid.New().String(), CarModel: CarModel{brand, model}}
//	}

type CarModel struct {
	Brand string
	Model string
}

type RentState uint

const (
	NotAvailable RentState = iota
	Available
)

type MaintenanceState uint

const (
	InMaintenance MaintenanceState = iota
	NotNeeded
	Needed
)

type Car struct {
	aggregate.ID
	VIN string
	CarModel
	RentState
	MaintenanceState
}

func New(model string, brand string) *Car {
	id, err := aggregate.NewID()
	if err != nil {
		panic(err)
	}
	return &Car{ID: id, VIN: uuid.New().String(), CarModel: CarModel{brand, model}}
}

func (c *Car) Register(car *Car) (aggregate.Events[Car], error) {
	return aggregate.NewEvents(&Arrived{Car: car}), nil
}

func (c *Car) ToCarV1() *carpark.Car {
	return &carpark.Car{
		ID:  c.ID,
		VIN: c.VIN,
		CarModel: carpark.CarModel{
			Brand: c.Brand,
			Model: c.Model,
		},
		RentState:        carpark.RentState(c.RentState),
		MaintenanceState: carpark.MaintenanceState(c.MaintenanceState),
	}
}
