package carpark

import "github.com/alekseev-bro/ddd/pkg/aggregate"

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

type MaintananceState uint

const (
	InMaintenance MaintananceState = iota
	NotNeeded
	Needed
)

type Car struct {
	aggregate.ID
	VIN string
	CarModel
	RentState
	MaintananceState
}

// func (c *Car) Reduce(event *gonvex.CoreEvent) error {

// 	switch event.Type {
// 	case CAR_CREATED:
// 		var car Car
// 		if err := json.Unmarshal(event.Payload, &car); err != nil {
// 			return fmt.Errorf("apply func: %w", err)
// 		}
// 		*c = car

// 	}

// 	return nil
// }
